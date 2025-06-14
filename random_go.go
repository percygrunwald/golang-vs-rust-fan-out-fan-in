package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var n int
	flag.IntVar(&n, "n", 50e6, "Number of random numbers to generate")
	flag.Parse()
	const N = 50e6 // Number of random numbers to generate

	start := time.Now()

	for i := 0; i < n; i++ {
		_ = rand.Intn(10) + 1 // Generate random number between 1 and 10
	}

	duration := time.Since(start)
	fmt.Printf("Go: Generated %d random numbers in %v\n", n, duration)
}
