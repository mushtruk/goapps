package concurrency

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestResultProcessor(t *testing.T) {
	inputs := []int{1, 2, 3, 4, 5}
	expectedSquares := []int{1, 4, 9, 16, 25}

	in := make(chan int, len(inputs))
	out := make(chan int, len(inputs))
	var wg sync.WaitGroup

	// Start the inputProducer to send values to the in channel
	go InputProducer(inputs, in)

	// Start the squareWorker to process the values from in and send to out
	wg.Add(1)
	go SquareWorker(in, out, &wg)

	// Close the out channel once all workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// Use ResultProcessor to collect the results from the out channel
	gotSquares := ResultProcessor(out)

	// Since the order does not matter, we have to sort both slices before comparing
	sort.Ints(expectedSquares)
	sort.Ints(gotSquares)

	if !reflect.DeepEqual(gotSquares, expectedSquares) {
		t.Errorf("ResultProcessor got %v, want %v", gotSquares, expectedSquares)
	}
}
