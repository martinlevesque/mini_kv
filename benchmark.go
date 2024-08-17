package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func safeAtoi(s string) int {
	i, err := strconv.Atoi(s)

	if err != nil {
		log.Println("Error converting string to int:", err)
		os.Exit(1)
	}

	return i
}

func main() {
	// CLI command arguments: go run benchmark.go <addr> <concurrency> <requests-per-second>

	// Get the address from the command line arguments
	addr := os.Args[1]
	concurrency := safeAtoi(os.Args[2])
	requestsPerSecond := safeAtoi(os.Args[3])

	log.Println("Starting benchmark to", addr)
	log.Println("Concurrency:", concurrency)
	log.Println("Requests per second:", requestsPerSecond)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done() // Signal completion of the goroutine
			benchmarkConnection(addr, requestsPerSecond)
		}()
		//go benchmarkConnection(addr, requestsPerSecond)
	}

	wg.Wait()
}

func benchmarkConnection(addr string, requestsPerSecond int) {
	log.Println("Benchmarking connection to", addr)
	// Create a connection to the server
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("Connected to", addr)
	time.Sleep(1 * time.Second)
}
