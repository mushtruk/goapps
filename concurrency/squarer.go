package concurrency

import (
	"sync"
)

// inputProducer sends integers to the input channel and closes the channel after.
func InputProducer(nums []int, in chan<- int) {
	// For each number in the slice, send it to the 'in' channel
	for _, number := range nums {
		in <- number
	}
	// Close the 'in' channel to indicate that's all the numbers
	close(in)
}

// squareWorker receives integers, squares them, and sends them to the out channel.
func SquareWorker(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	// For each integer in the 'in' channel
	// Calculate the square of the integer
	// Send the square to the 'out' channel

	for num := range in {
		out <- num * num
	}
	// Signal that this worker is done
	if wg != nil {
		wg.Done() // Notify the WaitGroup that this worker is done.
	}
}

// resultProcessor reads the squared values from the out channel and prints them.
func ResultProcessor(out <-chan int) []int {
	// For each square in the 'out' channel, print it
	var results []int
	for square := range out {
		results = append(results, square)
	}
	return results
}
