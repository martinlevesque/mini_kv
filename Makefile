
test:
	find . -name '*.go' | entr -r go test ./...

server:
	go run mini_kc.go

benchmark:
	go run benchmark.go localhost:8080 1000 10 10

tcp-client:
	nc localhost 8080

