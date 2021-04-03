package main

import (
	"net"
	"time"
)

func (b *backend) checkServer(index int) int {
	timeout := 3 * time.Second
	conn, err := net.DialTimeout("tcp", b.serverPool[index], timeout)

	if err != nil {
		return 0
	}

	defer conn.Close()

	if conn != nil {
		return 1
	}

	return 0
}

func (b *backend) healthCheck() {
	for {
		for i := 0; i < len(b.alive); i++ {
			b.alive[i] = b.checkServer(i)
		}
		time.Sleep(time.Second * 3)
	}
}
