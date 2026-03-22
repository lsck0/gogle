package main

import (
	"fmt"
	"net/http"
	"net/url"
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

func FetchWebPage(targetURL string) (WebPage, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return WebPage{}, fmt.Errorf("FetchWebPage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WebPage{}, fmt.Errorf("FetchWebPage: getting %q: %s", targetURL, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return WebPage{}, fmt.Errorf("FetchWebPage: parsing %q as HTML: %w", targetURL, err)
	}

	return parseWebPage(doc, resp.Request.URL), nil
}

func parseWebPage(doc *html.Node, url *url.URL) WebPage {
	return WebPage{
		URL:           url.String(),
		Title:         extractTitle(doc),
		Description:   extractDescription(doc),
		Content:       extractContent(doc),
		Links:         extractLinks(doc, url),
		RetrievalTime: time.Now(),
	}
}

func extractDescription(doc *html.Node) string {
	for n := range doc.Descendants() {
		if n.Type != html.ElementNode || n.DataAtom != atom.Meta {
			continue
		}

		found := false
		description := ""
		for _, a := range n.Attr {
			if a.Key == "property" && a.Val == "og:description" ||
				a.Key == "name" && a.Val == "description" {

				found = true
			} else if a.Key == "content" {
				description = a.Val
			}
		}
		if found {
			return description
		}
	}

	// not found
	return ""
}

func extractTitle(doc *html.Node) string {
	// try: open graph meta ("og:title")
	for n := range doc.Descendants() {
		if n.Type != html.ElementNode || n.DataAtom != atom.Meta {
			continue
		}

		found := false
		title := ""
		for _, a := range n.Attr {
			if a.Key == "property" && a.Val == "og:title" {
				found = true
			} else if a.Key == "content" {
				title = a.Val
			}
		}
		if found {
			return title
		}
	}

	// try: title tag
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Title && n.FirstChild != nil {
			return n.FirstChild.Data
		}
	}

	// no title found
	return ""
}

func extractContent(doc *html.Node) string {
	var sb strings.Builder

	// use recursion here so we can skip style and script element nodes and
	// all of their children, which Node.Descendents will not do.
	var inner func(n *html.Node)
	inner = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				sb.WriteString(text)
				sb.WriteByte(' ')
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode &&
				(c.DataAtom == atom.Style || c.DataAtom == atom.Script) {
				continue
			}
			inner(c)
		}
	}
	inner(doc)

	return strings.TrimSpace(sb.String()) // remove trailing whitespace
}

func extractLinks(doc *html.Node, url *url.URL) map[string]bool {
	links := make(map[string]bool)
	for n := range doc.Descendants() {
		if n.Type != html.ElementNode || n.DataAtom != atom.A {
			continue
		}
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			link, err := url.Parse(a.Val)
			if err != nil {
				continue // ignore error
			}
			normalizedLink, err := urlx.Normalize(link)
			if err != nil {
				continue // ignore error
			}

			if !links[normalizedLink] {
				links[normalizedLink] = true
			}
		}
	}
	return links
}
