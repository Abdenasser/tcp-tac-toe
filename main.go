package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Player struct {
	Connection net.Conn
	Symbol     string
	Score      int
}

// will only handle 2 players for now
var players = []Player{}
var lastPlayed = map[string]int{
	"playerIndex": 0,
}

func handleConnection(conn net.Conn, b map[int]string) {
	fmt.Printf("Client connected from %v \n", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	for {
		// output := printBoard(b)
		// conn.Write([]byte(output + "\n"))
		conn.Write([]byte("pick a position please \n"))
		ok := scanner.Scan()
		if !ok {
			break
		}
		play(scanner.Text(), b, players, conn)
	}
	removePlayer(conn)
	fmt.Printf("Client at %v disconnected.\n", conn.RemoteAddr().String())
	conn.Close()
}

func play(pos string, b map[int]string, ps []Player, conn net.Conn) {
	fmt.Println("> " + pos)
	for _, p := range ps {
		if lastPlayed["playerIndex"] == 0 {
			lastPlayed["playerIndex"] = 1
		} else {
			lastPlayed["playerIndex"] = 0
		}
		player := ps[lastPlayed["playerIndex"]]
		if p != player {
			position, _ := strconv.Atoi(pos)
			b[position] = player.Symbol
		}

		output := printBoard(b)
		p.Connection.Write([]byte(output + "\n"))
	}
}

func removePlayer(c net.Conn) {
	for i, p := range players {
		if p.Connection == c {
			players = append(players[:i], players[i+1:]...)
			break
		}
	}
}

func initBoard() map[int]string {
	board := make(map[int]string)
	for i := 0; i < 9; i++ {
		board[i] = fmt.Sprintf("%d", i)
	}
	return board
}

func printBoard(b map[int]string) string {
	output := ""
	for i := 0; i < 9; i += 3 {
		output += fmt.Sprintf("%v | %v | %v\n", b[i], b[i+1], b[i+2])
	}
	return output
}

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Printf("Listening on %s.\n", "localhost:8080")
	b := initBoard()
	defer listener.Close()

	for {
		conn, err := listener.Accept() // connect using telnet cmd: telnet localhost 8080
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		fmt.Printf("client connected from %v\n", conn.RemoteAddr().String())

		if len(players) < 2 {
			player := Player{
				Connection: conn,
				Symbol:     []string{"x", "o"}[len(players)],
				Score:      0,
			}
			players = append(players, player)
			go handleConnection(conn, b)
		}
	}
}
