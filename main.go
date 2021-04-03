package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type backend struct {
	serverPool []string
	alive      []int
	counter    int
}

var b backend
var stack_counter int

func bridge(wc io.WriteCloser, rd io.Reader) {
	defer wc.Close()
	io.Copy(wc, rd)
}

func handleConn(con net.Conn, addr string) error {
	be, err := net.Dial("tcp", addr)

	if err != nil {
		return errors.New("error: Connection cannot be handled!")
	}

	go bridge(con, be)
	go bridge(be, con)
	return nil
}

func (b *backend) selector() string {
	if stack_counter >= len(b.serverPool) {
		log.Fatal("All servers are crashed! Exiting...")
		os.Exit(2)
	}
	arrLength := len(b.serverPool)

	b.counter = b.counter % arrLength
	selection := b.counter
	b.counter++
	selected := b.serverPool[selection]

	if b.alive[selection] == 0 {
		stack_counter++
		selected = b.selector()
	}

	stack_counter = 0

	return selected
}

var (
	bind    = flag.String("bind", "", "The address to listen")
	servers = flag.String("backends", "", "backend addresses(must be separated by commas e.g '-backends 1.1.1.1:8080,2.2.2.2:8081')")
)

func main() {
	stack_counter = 0
	flag.Parse()
	if *bind == "" {
		log.Fatal("You must enter a valid address.")
	}

	if *servers == "" {
		log.Fatal("You must enter least one backend address.")
	}

	log.Printf("Backends: %s\n", *servers)

	b = backend{counter: 0, serverPool: strings.Split(*servers, ",")}

	b.alive = make([]int, len(b.serverPool))

	for i := 0; i < len(b.alive); i++ {
		b.alive[i] = 1
	}

	go b.healthCheck()

	ln, err := net.Listen("tcp", *bind)
	defer ln.Close()

	if err != nil {
		log.Fatal(errors.New("Address binding failed!"))
		os.Exit(1)
	}
	log.Printf("Listening %s\n", *bind)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go func() {
			err := handleConn(conn, b.selector())
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}()
	}
}
