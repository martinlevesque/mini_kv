# mini_kv

Minimalist (~ 400 lines) key-value store written in golang for the sake of learning go. 
The KV store is a TCP-based in-memory only and supports a small number of operations, namely GET, SET, DEL, KEYS, and EXPIRE.

## Usage

Boot up the server

```bash
make server
```

Connect to the server via TCP and run some commands. For example, using `nc`:

```bash
nc localhost 8080
SET key value
(nil)
GET key
value
```

## Benchmark

```bash
make benchmark
```

The benchmark was run using 2 servers over 1 router, where the benchmark would send significant amount of requests to the server.
See below the results:

Fixed number of connections set to 10 with variable number of requests per seconds per connection.

----------------------------------------------
| Requests per second | Throughput (requests per second sent-received) | Average delay |
|---------------------|-----------------------------------------------|---------------|
| 3                | 30.16                                     | 0.001056 sec.     |
| 30               | 295.00                                    | 0.001172 sec.     |
| 300              | 2454.51                                   | 0.000810 sec.     |
| 3000             | 21782.36                                  | 0.000249 sec.     |
| 30000            | 58929.76                                  | 0.000169 sec.     |
|---------------------|-----------------------------------------------|---------------|

Next we vary the number of connections with fixed number of requests per second per connection (10).


----------------------------------------------
| Num of Connections | Throughput (requests per second sent-received) | Average delay |
|---------------------|-----------------------------------------------|---------------|
| 3                | 29.90                                     | 0.001093 sec.     |
| 30               | 299.00                                    | 0.001302 sec.     |
| 300              | 2982.53                                   | 0.001865 sec.     |
| 3000             | 28074.55                                  | 0.002216 sec.     |
| 10000            | 71731.83                                  | 0.057918 sec.     |
|---------------------|-----------------------------------------------|---------------|
