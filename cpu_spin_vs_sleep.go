package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func tightLoopNoSleep(duration time.Duration) {
	fmt.Println("Starting tight loop (no sleep)...")
	start := time.Now()
	for time.Since(start) < duration {
		// Do nothing
	}
	fmt.Println("Tight loop (no sleep) finished.")
}

func tightLoopWithSleep(duration time.Duration, sleepDuration time.Duration) {
	fmt.Printf("Starting loop with %v sleep...\n", sleepDuration)
	start := time.Now()
	for time.Since(start) < duration {
		time.Sleep(sleepDuration)
	}
	fmt.Printf("Loop with %v sleep finished.\n", sleepDuration)
}

func main() {
	modeFlag := flag.String("m", "", "The mode to run in: 'sleep' or 'nosleep', required")
	runDurationFlag := flag.String("r", "", "Total length of time for the program to run (e.g. '5s'), required")
	sleepDurationFlag := flag.String("s", "", "Length of time to sleep (e.g. '5ms'), required in 'sleep' mode")
	flag.Parse()

	// Validate arguments
	if *modeFlag == "" || *runDurationFlag == "" {
		fmt.Println("Error: -m and -r arguments must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	runDuration, err := time.ParseDuration(*runDurationFlag)
	if err != nil {
		panic("Unable to parse run duration")
	}

	switch *modeFlag {
	case "nosleep":
		fmt.Printf("Running in mode '%v' for duration '%v'...\n", *modeFlag, runDuration)
		tightLoopNoSleep(runDuration)
	case "sleep":
		if *sleepDurationFlag == "" {
			fmt.Println("Error: -s argument must be provided if mode is 'sleep'.")
			flag.Usage()
			os.Exit(1)
		}
		sleepDuration, err := time.ParseDuration(*sleepDurationFlag)
		if err != nil {
			panic("Unable to parse sleep duration")
		}
		fmt.Printf("Running in mode '%v' for duration '%v' with sleep duration '%v'...\n",
			*modeFlag, runDuration, sleepDuration)
		tightLoopWithSleep(runDuration, sleepDuration)
	default:
		fmt.Println("Unknown mode:", os.Args[1])
		os.Exit(1)
	}
}
