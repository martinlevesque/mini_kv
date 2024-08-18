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

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done() // Signal completion of the goroutine
			benchmarkConnection(addr, concurrency, requestsPerSecond, duration)
		}()
	}

	wg.Wait()
}

func benchmarkConnection(addr string, concurrency int, requestsPerSecond int, durationSeconds int) {
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
	sumLatencyMs := 0

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
		currentDelay := int(endTimeTransmission.Sub(currentTime).Milliseconds())
		sumLatencyMs += currentDelay

		nbMicroSecondsInSecond := 1000000
		randomInterval := rand.Intn(nbMicroSecondsInSecond)
		multiplier := nbMicroSecondsInSecond/2 + randomInterval

		waitMicroseconds := int((1.0 / float64(requestsPerSecond)) * float64(multiplier))
		waitMicrosecondsWithoutDelay := waitMicroseconds - (currentDelay * 1000)

		if waitMicrosecondsWithoutDelay > 0 {
			time.Sleep(time.Duration(waitMicrosecondsWithoutDelay) * time.Microsecond)
		}

		if elapsed.Seconds() > float64(durationSeconds) {
			avgDelay := (float64(sumLatencyMs) / float64(cntSent)) / 1000
			fmt.Printf("%d %d %d %f\n", concurrency, requestsPerSecond, cntSent, avgDelay)
			break
		}
	}
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
