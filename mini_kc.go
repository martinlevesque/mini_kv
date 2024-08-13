package main

import (
	"bufio"
	"github.com/martinlevesque/mini-kc/kv"
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

	kvStore := kv.NewKVStore()

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Fatalf("Failed to accept connection: %s", err)
			continue
		}

		go handleConnection(conn, kvStore, kvStore.MutableCommandsChannel)
	}
}

func handleConnection(conn net.Conn, kvStore *kv.KVStore, commands_channel chan kv.KVOperation) {
	defer conn.Close()

	log.Printf("Accepted connection from %s", conn.RemoteAddr())

	// Read the request
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Printf("Failed to read request: %s", err)
			return
		}

		log.Printf("Received request: %s", line)

		commandResponse, err := kv.HandleCommand(line)

		if err != nil {
			log.Printf("Failed to handle command: %s", err)
			continue
		}

		if commandResponse.Action == kv.COMMAND_TERMINATE_CONN {
			log.Printf("Terminating connection")
			return
		}

		if commandResponse.Mutate {
			kvStore.MutableCommandsChannel <- commandResponse

			// Write the response
			_, err = conn.Write([]byte(<-commandResponse.ReplyCh + "\n"))
		} else {
			// Write the response
			log.Printf("immutable op")
			result, errOp := kvStore.ImmutableOperation(&commandResponse)

			if errOp != nil {
				log.Printf("Failed to do the immutable op, response: %s", errOp)
			} else {
				_, err = conn.Write([]byte(result + "\n"))

				if err != nil {
					log.Printf("Failed to write response: %s", err)
					return
				}
			}
		}

		if err != nil {
			log.Printf("Failed to write response: %s", err)
			return
		}
	}
}
