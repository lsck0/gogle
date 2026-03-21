package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goware/urlx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type WebPage struct {
	URL           string
	Title         string
	Description   string
	Content       string
	Links         map[string]bool
	RetrievalTime time.Time
}

func FetchWebPage(targetURL string) (res WebPage, err error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		err = fmt.Errorf("FetchWebPage: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("FetchWebPage: getting %q: %s", targetURL, resp.Status)
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		err = fmt.Errorf("FetchWebPage: parsing %q as HTML: %w", targetURL, err)
		return
	}

	res.URL = resp.Request.URL.String()
	res.Links = make(map[string]bool)
	res.RetrievalTime = time.Now()

	for n := range doc.Descendants() {
		// skip script and style
		if n.Type == html.ElementNode && (n.DataAtom == atom.Script || n.DataAtom == atom.Style) {
			continue
		}

		// meta data
		if n.Type == html.ElementNode && n.DataAtom == atom.Title {
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				var text = strings.TrimSpace(n.FirstChild.Data)
				if text != "" {
					res.Title = text
				}
			}
		}

		if n.Type == html.ElementNode && n.DataAtom == atom.Meta {
			var target = ""
			var content = ""
			for _, meta := range n.Attr {
				if meta.Key == "property" && meta.Val == "og:title" {
					target = "title"
				}

				if meta.Key == "name" && meta.Val == "description" {
					target = "description"
				}
				if meta.Key == "property" && meta.Val == "og:description" {
					target = "description"
				}

				if meta.Key == "content" {
					content = meta.Val
				}
			}

			switch target {
			case "title":
				res.Title = content
			case "description":
				res.Description = content
			}
		}

		// content
		if n.Type == html.TextNode {
			var text = strings.TrimSpace(n.Data)
			if text != "" {
				res.Content += text
			}
		}

		// links
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					var rawUrl = a.Val
					if !strings.HasPrefix(a.Val, "http://") && !strings.HasPrefix(a.Val, "https://") {
						rawUrl = res.URL + "/" + a.Val
					}
					var normalizedUrl, error = urlx.NormalizeString(rawUrl)
					if error != nil {
						continue
					}

					if !res.Links[normalizedUrl] {
						res.Links[normalizedUrl] = true
					}

					break
				}
			}
		}
	}

	return
}
