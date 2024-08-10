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

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Fatalf("Failed to accept connection: %s", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	// Read the request
	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalf("Failed to read request: %s", err)
		return
	}

	log.Printf("Received request: %s", line)

	HandleCommand(line)

}
