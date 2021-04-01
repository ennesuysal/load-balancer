package main

import (
	"errors"
	"io"
	"log"
	"net"
	"flag"
	"os"
	"strings"
)

type backend struct {
	serverPool []string
	counter int
}

var b backend

func bridge(wc io.WriteCloser, rd io.Reader){
	defer wc.Close()
	io.Copy(wc, rd)
}

func handleConn(con net.Conn, addr string) error{
	be, err := net.Dial("tcp", addr)

	if err != nil {
		return errors.New("error: Connection cannot be handled!")
	}

	go bridge(con, be)
	go bridge(be, con)
	return nil
}



func (b* backend)selector() string{
	arrLength := len(serverPool)

	b.counter = b.counter % arrLength
	selection := b.counter
	b.counter++

	return b.serverPool[selection]
}

var (
	bind = flag.String("bind", "", "The address to listen")
	servers = flag.String("backends", "", "backend addresses(must be separated by commas e.g '-backends 1.1.1.1, 2.2.2.2')")
)

func main(){
	flag.Parse()
	if *bind == "" {
		log.Fatal("You must enter a valid address.")
	}

	if *servers == "" {
		log.Fatal("You must enter least one backend address.")
	}

	log.Printf("Backends: %s\n", *servers)

	b = backend{counter: 0, serverPool: strings.Split(*servers, ",")}


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
		go func(){
			err := handleConn(conn, b.selector())
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}()
	}
}