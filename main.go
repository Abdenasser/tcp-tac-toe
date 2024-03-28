package main

import (
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Printf("Listening on %s.\n", "localhost:8080")

	for {
		conn, err := listener.Accept() // connect using telnet cmd: telnet localhost 8080
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}
		fmt.Printf("client connected from %v\n", conn.RemoteAddr().String())
	}
}
