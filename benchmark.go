//go:build benchmark
// +build benchmark

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type ConnectionResult struct {
	Concurrency       int
	RequestsPerSecond int
	Throughput        float64
	TotalRequests     int
	AvgLatency        float64
}

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
	// CLI command arguments: go run benchmark.go <addr> <concurrency> <requests-per-second> <duration>

	// Get the address from the command line arguments
	addr := os.Args[1]
	concurrency := safeAtoi(os.Args[2])
	requestsPerSecond := safeAtoi(os.Args[3])
	duration := safeAtoi(os.Args[4])

	results := make(chan ConnectionResult, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done() // Signal completion of the goroutine
			result := benchmarkConnection(addr, concurrency, requestsPerSecond, duration)

			results <- result
		}()
	}

	wg.Wait()
	close(results)

	// Calculate the average throughput and latency
	totalRequests := 0
	totalThroughput := 0.0
	totalLatency := 0.0

	for result := range results {
		totalRequests += result.TotalRequests
		totalThroughput += result.Throughput
		totalLatency += result.AvgLatency
	}

	avgThroughput := totalThroughput / float64(duration)
	avgLatency := totalLatency / float64(concurrency)

	fmt.Printf("%d %d %d %f %f\n", concurrency, requestsPerSecond, totalRequests, avgThroughput, avgLatency)
}

func benchmarkConnection(addr string, concurrency int, requestsPerSecond int, durationSeconds int) ConnectionResult {
	result := ConnectionResult{
		Concurrency:       concurrency,
		RequestsPerSecond: requestsPerSecond,
		Throughput:        0.0,
		TotalRequests:     0,
		AvgLatency:        0.0,
	}
	start := time.Now()

	// Create a connection to the server
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}

	defer conn.Close()
	reader := bufio.NewReader(conn)
	cntSent := 0
	sumLatencyMicrosecs := 0

	genCommands := map[int]func() string{
		0: genGet,
		1: genSet,
		2: genDel,
	}

	for {
		currentTime := time.Now()
		elapsed := currentTime.Sub(start)

		randomNumber := rand.Intn(3)
		funcToCall := genCommands[randomNumber]
		command := funcToCall()
		_, errWrite := conn.Write([]byte(command))

		if errWrite != nil {
			log.Println("Error writing to connection:", errWrite)
			break
		}

		_, errRead := reader.ReadString('\n')

		if errRead != nil {
			log.Println("Error reading from connection:", errRead)
			break
		}

		endTimeTransmission := time.Now()
		cntSent++
		currentDelay := int(endTimeTransmission.Sub(currentTime).Microseconds())
		sumLatencyMicrosecs += currentDelay

		waitMicroseconds := int(float64((1.0 / float64(requestsPerSecond))) * 1000000)
		waitMicrosecondsWithoutDelay := waitMicroseconds - (currentDelay)

		if waitMicrosecondsWithoutDelay > 0 {
			time.Sleep(time.Duration(waitMicrosecondsWithoutDelay) * time.Microsecond)
		}

		if elapsed.Seconds() > float64(durationSeconds) {
			break
		}
	}

	avgDelay := (float64(sumLatencyMicrosecs) / float64(cntSent)) / 1000000
	result.Throughput = float64(cntSent)
	result.TotalRequests = cntSent
	result.AvgLatency = avgDelay

	return result
}

func genGet() string {
	randomNumber := rand.Intn(100000)

	return fmt.Sprintf("GET key-%d\n", randomNumber)
}

func genSet() string {
	randomKey := rand.Intn(100000)
	randomValue := rand.Intn(100000)

	return fmt.Sprintf("SET key-%d \"%d\"\n", randomKey, randomValue)
}

func genDel() string {
	randomKey := rand.Intn(100000)

	return fmt.Sprintf("DEL key-%d\n", randomKey)
}
