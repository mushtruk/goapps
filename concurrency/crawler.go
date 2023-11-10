package concurrency

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func FetchUrl(url string) (body string, err error) {
	// Fetch the webpage
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and convert the body to string
	bodyBytes, err := io.ReadAll(resp.Body) // Corrected here

	if err != nil {
		return "", err
	}

	content := string(bodyBytes) // Corrected here

	return content, nil
}

func parseDom(n *html.Node, urls []string) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				urls = append(urls, a.Val)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		urls = parseDom(c, urls)
	}
	return urls
}

func ParseContent(content string) ([]string, error) {
	dom, err := html.Parse(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	urls = parseDom(dom, urls)

	return urls, nil
}
