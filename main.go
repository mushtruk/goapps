package main

import (
	"fmt"
	"sync"

	"github.com/mushtruk/goapps/concurrency/concurrency"
)

func main() {
	// Input slice of numbers
	nums := []int{1, 2, 3, 4, 5}

	// Create channels
	in := make(chan int, len(nums))
	out := make(chan int, len(nums))

	// Set up a WaitGroup to wait for all squareWorkers to finish
	var wg sync.WaitGroup

	// Start the inputProducer as a goroutine
	go concurrency.InputProducer(nums, in)

	// Define the number of workers
	numWorkers := 3
	wg.Add(numWorkers)

	// Start squareWorkers
	for i := 0; i < numWorkers; i++ {
		go concurrency.SquareWorker(in, out, &wg)
	}

	// Start a goroutine to close 'out' once all workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// Read from out and print each result
	for sqr := range out {
		fmt.Println(sqr)
	}
}
