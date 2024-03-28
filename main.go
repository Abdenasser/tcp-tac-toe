package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Player struct {
	Index      int
	Connection net.Conn
	Symbol     string
	Score      int
}

// will only handle 2 players for now
var players = []Player{}
var currentPlayer Player

const NoWinner = "none"

var colors = map[string][2]string{
	"orange": {"\x1b[34m", "\x1b[0m"},
	"cyan":   {"\x1b[36m", "\x1b[0m"},
}

func colorize(text string, color string) string {
	return fmt.Sprintf("%s%s%s", colors[color][0], text, colors[color][1])
}

func closeConnection(c net.Conn) {
	for i, p := range players {
		if p.Connection == c {
			players = append(players[:i], players[i+1:]...)
			break
		}
	}
	c.Close()
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		play(scanner.Text(), conn)
	}
	closeConnection(conn)

}

func switchCurrentPlayer() {
	if currentPlayer.Index == 0 {
		currentPlayer = players[1]
	} else {
		currentPlayer = players[0]
	}
}

func getScore() string {
	if len(players) == 2 {
		return fmt.Sprintf("score: %v:%v - %v:%v", players[0].Symbol, players[0].Score, players[1].Symbol, players[1].Score)
	}
	return "waiting for an oponent to join"
}

func dispatchBoard() {
	for _, p := range players {
		output := printBoard(originalboard)
		score := getScore()
		p.Connection.Write([]byte(output + "\n"))
		p.Connection.Write([]byte(score + "\n"))
		if p.Connection == currentPlayer.Connection && len(players) == 2 {
			p.Connection.Write([]byte("your turn \n"))
		} else {
			p.Connection.Write([]byte("oponent's turn \n"))
		}
	}
}

func isFreePos(pos int) bool {
	if originalboard[pos] == colorize("x", "orange") || originalboard[pos] == colorize("o", "cyan") {
		return false
	}
	return true
}

func getWinner(ps []Player) string {
	winCombos := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // verticals
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // horizontals
		{0, 4, 8}, {2, 4, 6}, // diagonals
	}

	for _, combo := range winCombos {
		if originalboard[combo[0]] == originalboard[combo[1]] && originalboard[combo[1]] == originalboard[combo[2]] {
			if originalboard[combo[0]] == ps[0].Symbol {
				ps[0].Score += 1
			}
			if originalboard[combo[0]] == ps[1].Symbol {
				ps[1].Score += 1
			}
			return originalboard[combo[0]]
		}
	}

	return NoWinner
}

func isFull(b map[int]string) bool {
	for _, v := range b {
		if v != colorize("x", "orange") && v != colorize("o", "cyan") {
			return false
		}
	}
	return true
}

func shouldReset(w string) bool {
	if w != NoWinner || (w == NoWinner && isFull(originalboard)) {
		return true
	}
	return false
}

func play(pos string, c net.Conn) {
	fmt.Println("> " + pos)
	position, _ := strconv.Atoi(pos)

	if c == currentPlayer.Connection && len(players) == 2 && isFreePos(position) {
		originalboard[position] = currentPlayer.Symbol
		switchCurrentPlayer()
	}
	winner := getWinner(players)
	if shouldReset(winner) {
		originalboard = initBoard()
	}
	dispatchBoard()
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

var originalboard = initBoard()

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Println("Listening on localhost:8080.")

	defer listener.Close()

	for {
		conn, _ := listener.Accept() // connect using telnet cmd: telnet localhost 8080
		fmt.Printf("client connected from %v\n", conn.RemoteAddr().String())

		if len(players) < 2 {
			player := Player{
				Index:      len(players),
				Connection: conn,
				Symbol:     []string{colorize("x", "orange"), colorize("o", "cyan")}[len(players)],
				Score:      0,
			}
			players = append(players, player)
			currentPlayer = player
			go handleConnection(conn)
		}
	}
}
