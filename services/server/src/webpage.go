package main

import (
	"github.com/goware/urlx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strings"
	"time"
)

type WebPage struct {
	Url           string
	Title         string
	Description   string
	Content       string
	Links         map[string]bool
	RetrievalTime time.Time
}

func ReadWebPage(targetUrl string) (result WebPage, ok bool) {
	response, err := http.Get(targetUrl)
	if err != nil {
		return
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return
	}

	result.Url = response.Request.URL.String()
	result.Links = make(map[string]bool)
	result.RetrievalTime = time.Now()

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
					result.Title = text
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
				result.Title = content
			case "description":
				result.Description = content
			}
		}

		// content
		if n.Type == html.TextNode {
			var text = strings.TrimSpace(n.Data)
			if text != "" {
				result.Content += text
			}
		}

		// links
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					var rawUrl = a.Val
					if !strings.HasPrefix(a.Val, "http://") && !strings.HasPrefix(a.Val, "https://") {
						rawUrl = result.Url + "/" + a.Val
					}
					var normalizedUrl, error = urlx.NormalizeString(rawUrl)
					if error != nil {
						continue
					}

					if !result.Links[normalizedUrl] {
						result.Links[normalizedUrl] = true
					}

					break
				}
			}
		}
	}

	ok = true
	return
}
