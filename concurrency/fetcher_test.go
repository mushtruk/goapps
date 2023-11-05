package concurrency_test

import (
	"sync"
	"testing"

	"github.com/mushtruk/goapps/concurrency/concurrency"
)

func TestFetch(t *testing.T) {
	urls := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		// Non-existent site for error checking
		"https://a-very-unlikely-domain-name.xyz",
	}

	var wg sync.WaitGroup
	resultsCh := make(chan concurrency.Result, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go concurrency.Fetch(url, &wg, resultsCh)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect all results
	for result := range resultsCh {
		if result.Err != nil {
			t.Logf("Expected error for URL %s: %v", result.URL, result.Err)
		} else if result.Length == 0 {
			t.Errorf("Expected non-zero content length for %s", result.URL)
		} else {
			t.Logf("Got response from %s, length: %d", result.URL, result.Length)
		}
	}
}
