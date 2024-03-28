package main

import (
	"fmt"
	"net"
)

var playerSymbols = []string{"x", "o"}
var connections = []net.Conn{}

type Player struct {
	Connection net.Conn
	Symbol     string
	Score      int
}

// will only handle 2 players for now
var players = []Player{}

func popSymbol(s *[]string) string {
	pop := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return pop
}

func createPlayer(c net.Conn) Player {
	player := Player{
		Connection: c,
		Symbol:     popSymbol(&playerSymbols),
		Score:      0,
	}
	return player
}

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Printf("Listening on %s.\n", "localhost:8080")

	for {
		conn, err := listener.Accept() // connect using telnet cmd: telnet localhost 8080
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}
		fmt.Printf("client connected from %v\n", conn.RemoteAddr().String())
		connections = append(connections, conn)
		if len(players) < 2 {
			players = append(players, createPlayer(conn))
		}
		fmt.Printf("clients %v, players %v, symbols %v\n", len(connections), len(players), len(playerSymbols))
	}
}
