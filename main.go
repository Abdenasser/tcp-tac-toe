package main

import (
	"fmt"
	"net"
	"time"
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

func handleConnection(conn net.Conn) {
	fmt.Printf("Client connected from %v \n", conn.RemoteAddr().String())
	time.Sleep(5 * time.Second)
	closeConnection(conn)
}

func closeConnection(conn net.Conn) {
	removeConn(conn)
	removePlayer(conn)
	fmt.Printf("Client at %v disconnected.\n", conn.RemoteAddr().String())
	conn.Close()
}

func removeConn(c net.Conn) {
	for i, conn := range connections {
		if conn == c {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

func removePlayer(c net.Conn) {
	for i, p := range players {
		if p.Connection == c {
			playerSymbols = append(playerSymbols, p.Symbol)
			players = append(players[:i], players[i+1:]...)
			break
		}
	}
}

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Printf("Listening on %s.\n", "localhost:8080")

	defer listener.Close()

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
		// accept inputs from connected clients in a separate goroutine
		go handleConnection(conn)
		fmt.Printf("clients %v, players %v, symbols %v\n", len(connections), len(players), len(playerSymbols))
	}
}
