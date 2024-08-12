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
	commands_channel := make(chan kv.KVOperation)

	go func(kvOperation <-chan kv.KVOperation) {
		for {
			// Wait for a command from the channel
			currentKvOperation := <-kvOperation

			log.Println("opppp Received command: ", currentKvOperation)

			if currentKvOperation.Action == kv.COMMAND_SET_KEY {
				log.Println("set value")
				kvStore.Set(currentKvOperation.KeyName, currentKvOperation.Value)
			}

			currentKvOperation.ReplyCh <- "Hello from the server!"
		}
	}(commands_channel)

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Fatalf("Failed to accept connection: %s", err)
			continue
		}

		go handleConnection(conn, kvStore, commands_channel)
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

		replyCh := make(chan string)
		commandResponse.ReplyCh = replyCh

		if commandResponse.Action == kv.COMMAND_TERMINATE_CONN {
			log.Printf("Terminating connection")
			return
		}

		if commandResponse.Mutate {
			commands_channel <- commandResponse

			// Write the response
			_, err = conn.Write([]byte(<-replyCh + "\n"))
		} else {
			// Write the response
			// todo handle get
			keyValue, errGet := kvStore.Get(commandResponse.KeyName)

			if errGet != nil {
				_, err = conn.Write([]byte("(nil)\n"))
			} else {
				_, err = conn.Write([]byte("\"" + keyValue + "\"\n"))
			}
		}

		if err != nil {
			log.Printf("Failed to write response: %s", err)
			return
		}
	}
}
