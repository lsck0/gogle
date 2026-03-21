package collection

import (
	"log"
	"sync"
)

// this is a wrapper around a channel and a buffer string
// buffer slice add the capability of adding burst of data in the Queue
// strem add the capability of making the Queue a channel type consumable
type URLQueue struct {
	mu                 sync.RWMutex
	buffer             []string
	intermittentBuffer []string
	stream             chan string
	seen               map[string]struct{}
	signal             chan bool // this will have value if buffer has pending data to be sent on stream
}

func NewURLQueue() *URLQueue {
	obj := &URLQueue{
		mu:     sync.RWMutex{},
		buffer: make([]string, 0),
		signal: make(chan bool, 1), // to avoid blocking at push during the stream is not able to cope up backpressure
		stream: make(chan string),
		seen:   make(map[string]struct{}),
	}
	go obj.workerFromBufferToStream()
	return obj
}

func (r *URLQueue) GetStream() <-chan string {
	// read only stream just casting of chan string to <-chan string
	return r.stream
}

func (r *URLQueue) PushUrls(urls ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	someNewDataAdded := false
	for _, url := range urls {
		if _, ok := r.seen[url]; !ok {
			r.buffer = append(r.buffer, url)
			r.seen[url] = struct{}{}
			someNewDataAdded = true
		}
	}

	if someNewDataAdded {
		log.Println("DEBUG : urls pending to be crawled", "count", len(r.buffer)+len(r.intermittentBuffer))

		select {
		case r.signal <- true:
		default:
			// no issue as already a pending signal is their
		}
	}
}

func (r *URLQueue) workerFromBufferToStream() {
	// a continuous loop that waits for a signal to copy new buffer data to stream
	for {
		<-r.signal

		r.mu.Lock()
		r.intermittentBuffer = r.buffer
		r.buffer = make([]string, 0)
		r.mu.Unlock()

		for _, url := range r.intermittentBuffer {
			r.stream <- url
		}
	}
}
