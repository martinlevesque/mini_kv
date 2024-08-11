package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	log.Println("Starting Mini-KC")

	// Listen on TCP port 8080
	server, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
		os.Exit(1)
	}

	defer server.Close()

	commands_channel := make(chan KVOperation)

	go func(kvOperation <-chan KVOperation) {
		for {
			// Wait for a command from the channel
			currentKvOperation := <-kvOperation

			log.Println("opppp Received command: ", currentKvOperation)

			currentKvOperation.replyCh <- "Hello from the server!"
		}
	}(commands_channel)

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Fatalf("Failed to accept connection: %s", err)
			continue
		}

		go handleConnection(conn, commands_channel)
	}
}

func handleConnection(conn net.Conn, commands_channel chan KVOperation) {
	defer conn.Close()

	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	// Read the request
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("Failed to read request: %s", err)
			return
		}

		log.Printf("Received request: %s", line)

		commandResponse, err := HandleCommand(line)

		if err != nil {
			log.Printf("Failed to handle command: %s", err)
			continue
		}

		replyCh := make(chan string)
		commandResponse.replyCh = replyCh

		if commandResponse.Action == COMMAND_TERMINATE_CONN {
			log.Printf("Terminating connection")
			return
		}

		commands_channel <- commandResponse

		// Write the response
		_, err = conn.Write([]byte(<-replyCh + "\n"))

		if err != nil {
			log.Printf("Failed to write response: %s", err)
			return
		}
	}
}
