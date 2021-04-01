## A simple Load Balancer

### Building

Run it in same directory with project root.

	go build .

### Usage

	./load-balancer -bind 0.0.0.0:8083 -backends 127.0.0.1:8081,127.0.0.1:8082