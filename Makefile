
test:
	find . -name '*.go' | entr -r go test ./...

server:
	go run mini_kc.go

benchmark:
	# go run benchmark.go  <addr> <concurrency> <requests-per-second> <duration>
	echo "Concurrency,Requests per second,Throughput per second,Total requests,Avg delay" > benchmark-rps.csv
	go run benchmark.go localhost:8080 10 3 60 >> benchmark-rps.csv
	go run benchmark.go localhost:8080 10 30 60 >> benchmark-rps.csv
	go run benchmark.go localhost:8080 10 300 60 >> benchmark-rps.csv
	go run benchmark.go localhost:8080 10 3000 60 >> benchmark-rps.csv
	go run benchmark.go localhost:8080 10 30000 60 >> benchmark-rps.csv

	echo "Concurrency,Requests per second,Throughput per second,Total requests,Avg delay" > benchmark-conns.csv
	go run benchmark.go localhost:8080 3 10 60 >> benchmark-conns.csv
	go run benchmark.go localhost:8080 30 10 60 >> benchmark-conns.csv
	go run benchmark.go localhost:8080 300 10 60 >> benchmark-conns.csv
	go run benchmark.go localhost:8080 3000 10 60 >> benchmark-conns.csv
	go run benchmark.go localhost:8080 30000 10 60 >> benchmark-conns.csv

tcp-client:
	nc localhost 8080

