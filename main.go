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

type Game struct {
	Board         map[int]string
	Players       []Player
	CurrentPlayer Player
	Watchers      []net.Conn
}

func (g *Game) printBoard() string {
	output := ""
	for i := 0; i < 9; i += 3 {
		output += fmt.Sprintf("%v | %v | %v\n", g.Board[i], g.Board[i+1], g.Board[i+2])
	}
	return output
}

func (g *Game) isFullBoard() bool {
	for _, v := range g.Board {
		if v != colorize("x", "orange") && v != colorize("o", "cyan") {
			return false
		}
	}
	return true
}

func (g *Game) shouldResetBoard() bool {
	if g.getWinner() != "none" || (g.getWinner() == "none" && g.isFullBoard()) {
		return true
	}
	return false
}

func (g *Game) canPlayTurn(p Player) bool {
	if p.Connection == g.CurrentPlayer.Connection && len(g.Players) == 2 {
		return true
	}
	return false
}

func (g *Game) resetBoard() {
	for i := range 9 {
		g.Board[i] = fmt.Sprintf("%d", i)
	}
}

func (g *Game) switchCurrentPlayer() {
	if g.CurrentPlayer.Index == 0 {
		g.CurrentPlayer = g.Players[1]
	} else {
		g.CurrentPlayer = g.Players[0]
	}
}

func (g *Game) getScore() string {
	if len(g.Players) == 2 {
		return fmt.Sprintf("score: %v:%v - %v:%v", g.Players[0].Symbol, g.Players[0].Score, g.Players[1].Symbol, g.Players[1].Score)
	}
	return "waiting for an oponent to join"
}

func (g *Game) isFreePos(pos int) bool {
	if g.Board[pos] == colorize("x", "orange") || g.Board[pos] == colorize("o", "cyan") {
		return false
	}
	return true
}

func (g *Game) playTurn(p Player, pos int) {
	g.Board[pos] = p.Symbol
}

func (g *Game) getWinner() string {
	winCombos := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // verticals
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // horizontals
		{0, 4, 8}, {2, 4, 6}, // diagonals
	}

	for _, combo := range winCombos {
		if g.Board[combo[0]] == g.Board[combo[1]] && g.Board[combo[1]] == g.Board[combo[2]] {
			if g.Board[combo[0]] == g.Players[0].Symbol {
				g.Players[0].Score += 1
			} else if g.Board[combo[0]] == g.Players[1].Symbol {
				g.Players[1].Score += 1
			}
			return g.Board[combo[0]]
		}
	}

	return "none"
}

func (g *Game) dispatch() {
	for _, p := range g.Players {
		if len(g.Players) != 2 {
			p.Connection.Write([]byte("wait for an oponent to join" + "\n"))
			return
		}
		p.Connection.Write([]byte(g.printBoard() + "\n" + g.getScore() + "\n"))
		if p.Connection == g.CurrentPlayer.Connection {
			p.Connection.Write([]byte("your turn \n"))
			continue
		}
		p.Connection.Write([]byte("oponent's turn \n"))
	}
}

var colors = map[string][2]string{
	"orange": {"\x1b[34m", "\x1b[0m"},
	"cyan":   {"\x1b[36m", "\x1b[0m"},
}

func colorize(text string, color string) string {
	return fmt.Sprintf("%s%s%s", colors[color][0], text, colors[color][1])
}

func handleGameConnection(g *Game, p Player) {
	scanner := bufio.NewScanner(p.Connection)

	defer handlePlayerQuit(g, p)

	for {
		ok := scanner.Scan()
		if !ok {
			break
		}
		pos, _ := strconv.Atoi(scanner.Text())
		handlePlayerPosition(pos, g, p)
	}

}

func handlePlayerQuit(g *Game, p Player) {
	for i, p := range g.Players {
		if p.Connection == p.Connection {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}
	p.Connection.Close()
}

func handlePlayerPosition(pos int, g *Game, p Player) {
	if g.canPlayTurn(p) && g.isFreePos(pos) {
		g.playTurn(p, pos)
		g.switchCurrentPlayer()
	}

	if g.shouldResetBoard() {
		g.resetBoard()
	}

	g.dispatch()
}

func main() {
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Println("Listening on localhost:8080.")

	var game = Game{
		Board: make(map[int]string, 9),
	}

	game.resetBoard()

	defer listener.Close()

	for {
		conn, _ := listener.Accept()

		if len(game.Players) == 2 {
			fmt.Printf("rejecting client connected from %v\n", conn.RemoteAddr().String())
			conn.Write([]byte("Game is full, try again later!" + "\n"))
			conn.Close() // only accepting 2 player
			continue
		}

		fmt.Printf("client connected from %v\n", conn.RemoteAddr().String())
		player := Player{
			Index:      len(game.Players),
			Connection: conn,
			Symbol:     []string{colorize("x", "orange"), colorize("o", "cyan")}[len(game.Players)],
			Score:      0,
		}
		game.Players = append(game.Players, player)
		game.CurrentPlayer = player

		if len(game.Players) == 2 {
			game.Players[0].Connection.Write([]byte("Oponent have joined" + "\n"))
			game.dispatch()
		}

		go handleGameConnection(&game, player)
	}
}
