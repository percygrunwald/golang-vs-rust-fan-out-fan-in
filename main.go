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
	numValues := flag.Int("n", 0, "Number of random integers to produce")
	numConsumers := flag.Int("w", 0, "Number of consumer goroutines")
	flag.Parse()

	// Validate arguments
	if *numValues <= 0 || *numConsumers <= 0 {
		fmt.Println("Error: Both -n and -w arguments must be provided and greater than 0.")
		flag.Usage()
		os.Exit(1)
	}

	const bufferSize = 1e3

	// Channels
	produceChan := make(chan int, bufferSize)
	squareChan := make(chan int, bufferSize)

	// WaitGroups
	var wgConsumers sync.WaitGroup
	var wgFanIn sync.WaitGroup

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
		for i := 0; i < *numValues; i++ {
			num := getNumber(i) // Get number from the list
			produceChan <- num
		}
		close(produceChan)
	}()

	// Consumer goroutines
	wgConsumers.Add(*numConsumers)
	for i := 0; i < *numConsumers; i++ {
		go func() {
			defer wgConsumers.Done()
			for num := range produceChan {
				squareChan <- num * num
			}
		}()
	}

	// Fan-in goroutine
	wgFanIn.Add(1)
	go func() {
		defer wgFanIn.Done()
		var sum int
		for square := range squareChan {
			sum += square
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
