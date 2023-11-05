package concurrency

import (
	"io"
	"net/http"
	"sync"
)

// Result is used to hold the URL and the length of its content
type Result struct {
	URL    string
	Length int
	Err    error
}

// Fetch fetches the content from the given URL and sends a Result to the provided channel
func Fetch(url string, wg *sync.WaitGroup, ch chan<- Result) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		ch <- Result{URL: url, Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- Result{URL: url, Err: err}
		return
	}

	ch <- Result{URL: url, Length: len(body)}
}
