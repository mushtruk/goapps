package concurrency_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mushtruk/goapps/concurrency/concurrency"
)

const content = `<html>
<head><title>Test Page</title></head>
<body>
    <a href="/link1">Link 1</a>
    <a href="/link2">Link 2</a>
</body>
</html>`

// startTestServer returns a new mock HTTP server that responds with a simple HTML page containing links.
func newTestServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, content)
	})
	return httptest.NewServer(handler) // This starts and returns a new server
}

func TestFetchUrl(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	content, err := concurrency.FetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("FetchUrl() returned an error: %v", err)
	}

	if content != content {
		t.Errorf("FetchUrl() got = %s, want %s", content, content)
	}
}

func TestParseContent(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	_content, _ := concurrency.FetchUrl(ts.URL)

	urls, err := concurrency.ParseContent(_content)

	if err != nil {
		t.Fatalf("Failed to fetch content from test server: %v", err)
	}

	expectedURLs := []string{"/link1", "/link2"}

	if len(urls) != len(expectedURLs) {
		t.Errorf("ParseContent() returned %d URLs, want %d", len(urls), len(expectedURLs))
	}

	for i, url := range urls {
		if url != expectedURLs[i] {
			t.Errorf("ParseContent() URLs[%d] = %s, want %s", i, url, expectedURLs[i])
		}
	}

}
