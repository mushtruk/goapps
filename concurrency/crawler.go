package concurrency

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Queue struct {
	items   []string
	visited map[string]bool
}

func NewQueue() *Queue {
	return &Queue{
		items:   make([]string, 0),
		visited: make(map[string]bool),
	}
}

func (q *Queue) Add(url string) {
	if !q.visited[url] {
		q.items = append(q.items, url)
	}
}

func (q *Queue) Next() string {
	if len(q.items) == 0 {
		return ""
	}

	url := q.items[0]
	q.items = q.items[:1]

	return url
}

func (q *Queue) MarkVisited(url string) {
	q.visited[url] = true
}

func (q *Queue) IsVisited(url string) bool {
	_, visited := q.visited[url]
	return visited
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue) Size() int {
	return len(q.items)
}

func FetchUrl(url string) (body string, err error) {
	// Fetch the webpage
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and convert the body to string
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	content := string(bodyBytes)

	return content, nil
}

func parseDom(n *html.Node, urls []string, base *url.URL) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				href, err := url.Parse(a.Val)
				if err != nil {
					continue // Handle or log the error as per your requirement
				}
				resolvedURL := base.ResolveReference(href).String()
				urls = append(urls, resolvedURL)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		urls = parseDom(c, urls, base)
	}
	return urls
}

func ParseContent(content string, baseURL string) ([]string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err // Handle error if base URL is invalid
	}

	dom, err := html.Parse(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	urls = parseDom(dom, urls, base)

	return urls, nil
}

func CrawlNextURL(q *Queue, base string) {

	next := q.Next()

	u, err := url.ParseRequestURI(next)

	if err != nil {
		return
	}

	content, err := FetchUrl(u.String())

	if err != nil {
		return
	}

	extractedUrls, err := ParseContent(content, base)

	if err != nil {
		return
	}

	for _, extractedUrl := range extractedUrls {
		if !q.IsVisited(extractedUrl) {
			q.Add(extractedUrl)
		}
	}
	q.MarkVisited(next)
}
