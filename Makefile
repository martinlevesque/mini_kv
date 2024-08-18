
test:
	find . -name '*.go' | entr -r go test ./...

server:
	go run mini_kc.go

benchmark:
	# go run benchmark.go <host> <num_requests> <num_concurrent> <num_requests_per_second> <duration>
	go run benchmark.go localhost:8080 10 10 10

tcp-client:
	nc localhost 8080

