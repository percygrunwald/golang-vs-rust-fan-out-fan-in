package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

func main() {
	/*
		Fan-Out and Fan-In Diagram:

		Producer Goroutine
					┌───────────────┐
					│ Random Numbers│
					│   Generator   │
					└───────┬───────┘
									│
									▼
					┌────────────────────┐
					│   produceChan      │
					│ (Buffered Channel) │
					└─────────┬──────────┘
										│
										▼
					┌───────────────────────────────────────────────────────┐
					│                                                       │
					│                Consumer Goroutines                    │
					│                                                       │
					│ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ │
					│ │ Square Numbers│ │ Square Numbers│ │ Square Numbers│ │
					│ │   Calculator  │ │   Calculator  │ │   Calculator  │ │
					│ └───────┬───────┘ └───────┬───────┘ └───────┬───────┘ │
					│         │                 │                 │         │
					└─────────▼─────────────────▼─────────────────▼─────────┘
										│
										▼
					┌────────────────────┐
					│   squareChan       │
					│ (Buffered Channel) │
					└─────────┬──────────┘
										│
										▼
					┌─────────────────┐
					│ Fan-In Goroutine│
					│   Sum Calculator│
					└───────┬─────────┘
									│
									▼
					┌────────────────┐
					│ Final Sum      │
					│   Printed      │
					└────────────────┘
	*/

	// Command line arguments
	batchSize := flag.Int("b", 0, "Specify the batch size")
	numValues := flag.Int("n", 0, "Number of random integers to produce")
	numConsumers := flag.Int("w", 0, "Number of consumer goroutines")
	flag.Parse()

	// Validate arguments
	if *numValues <= 0 || *numConsumers <= 0 || *batchSize <= 0 {
		fmt.Println("Error: -b, -n and -w arguments must be provided and greater than 0.")
		flag.Usage()
		os.Exit(1)
	}

	const bufferSize = 1e3

	// Channels
	produceChan := make(chan []int, bufferSize)
	squareChan := make(chan []int, bufferSize)

	// WaitGroups
	var wgConsumers sync.WaitGroup
	var wgFanIn sync.WaitGroup
	var mu sync.Mutex

	// Predefined list of random numbers generated using rand.Intn
	randomNumbers := make([]int, 10)
	for i := range randomNumbers {
		randomNumbers[i] = rand.Intn(10) + 1 // Random integer between 1 and 10
	}
	listLength := len(randomNumbers)

	// Function to cycle through the list
	getNumber := func(index int) int {
		return randomNumbers[index%listLength]
	}

	// Producer goroutine
	go func() {
		for i := 0; i < *numValues; i += *batchSize {
			batch := make([]int, 0, *batchSize)
			for j := 0; j < *batchSize && i+j < *numValues; j++ {
				num := getNumber(i + j) // Get number from the list
				mu.Lock()
				batch = append(batch, num)
				mu.Unlock()
			}
			produceChan <- batch
		}
		close(produceChan)
	}()

	// Consumer goroutines
	wgConsumers.Add(*numConsumers)
	for i := 0; i < *numConsumers; i++ {
		go func() {
			defer wgConsumers.Done()
			for batch := range produceChan {
				squaredBatch := make([]int, len(batch))
				for j, num := range batch {
					squaredBatch[j] = num * num
				}
				squareChan <- squaredBatch
			}
		}()
	}

	// Fan-in goroutine
	wgFanIn.Add(1)
	go func() {
		defer wgFanIn.Done()
		var sum int
		for batch := range squareChan {
			for _, square := range batch {
				sum += square
			}
		}
		fmt.Println("Final Sum:", sum)
	}()

	// Wait for consumers to finish and close squareChan
	go func() {
		wgConsumers.Wait()
		close(squareChan)
	}()

	// Wait for fan-in goroutine to finish
	wgFanIn.Wait()
	fmt.Println("Finished main thread.")
}
