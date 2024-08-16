
test:
	find . -name '*.go' | entr -r go test ./...

server:
	go run .

tcp-client:
	nc localhost 8080

