package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const xSymbol = "\x1b[34mX\x1b[0m"
const oSymbol = "\x1b[36mO\x1b[0m"

const noWinner = "none"

var winCombinations = [][]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // verticals
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // horizontals
	{0, 4, 8}, {2, 4, 6}, // diagonals
}

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

func (g *Game) String() string {
	if !g.isFullPlaces() {
		return "waiting for an oponent to join \n"
	}
	var (
		board string
		score string
	)
	for i := 0; i < 9; i += 3 {
		board += fmt.Sprintf("%v | %v | %v\n", g.Board[i], g.Board[i+1], g.Board[i+2])
	}
	for _, p := range g.Players {
		score += fmt.Sprintf("%v:%v ", p.Symbol, p.Score)

	}
	return fmt.Sprintf("%v \nscore: %v\n", board, score)
}

func (g *Game) isFullBoard() bool {
	for _, v := range g.Board {
		if v != xSymbol && v != oSymbol {
			return false
		}
	}
	return true
}

func (g *Game) isFullPlaces() bool {
	return len(g.Players) == 2
}

func (g *Game) shouldResetBoard() bool {
	return g.getWinnerSymbol() != noWinner || g.isFullBoard()
}

func (g *Game) canPlayTurn(p Player) bool {
	if p.Connection == g.CurrentPlayer.Connection && g.isFullPlaces() {
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

func (g *Game) isFreePos(pos int) bool {
	if g.Board[pos] == xSymbol || g.Board[pos] == oSymbol || pos > 8 || pos < 0 {
		return false
	}
	return true
}

func (g *Game) playTurn(p Player, pos int) {
	g.Board[pos] = p.Symbol
}

func (g *Game) getWinnerSymbol() string {
	for _, c := range winCombinations {
		if g.Board[c[0]] == g.Board[c[1]] && g.Board[c[1]] == g.Board[c[2]] {
			return g.Board[c[0]]
		}
	}
	return noWinner
}

func (g *Game) getPlayerWithSymbol(s string) *Player {
	if g.Players[0].Symbol == s {
		return &g.Players[0]
	}
	return &g.Players[1]
}

func (p *Player) incrementScore() {
	p.Score += 1
}

func (g *Game) dispatchNextTurn() {
	if !g.isFullPlaces() {
		return
	}
	for _, p := range g.Players {
		if p.Connection != g.CurrentPlayer.Connection {
			p.Connection.Write([]byte("oponent's turn! \n"))
			continue
		}
		p.Connection.Write([]byte("your turn: \n"))
	}
}

func (g *Game) dispatchGame() {
	for _, p := range g.Players {
		p.Connection.Write([]byte(g.String()))
	}
}

func handleGameConnection(g *Game, p Player) {
	scanner := bufio.NewScanner(p.Connection)

	defer handlePlayerQuit(g, p)

	for {
		if ok := scanner.Scan(); !ok {
			break
		}
		pos, _ := strconv.Atoi(scanner.Text())
		handlePlayerPosition(pos, g, p)
	}

}

func handlePlayerQuit(g *Game, pl Player) {
	for i, p := range g.Players {
		if p.Connection == pl.Connection {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			break
		}
	}
	pl.Connection.Close()
}

func handlePlayerPosition(pos int, g *Game, p Player) {
	if g.canPlayTurn(p) && g.isFreePos(pos) {
		g.playTurn(p, pos)
		g.switchCurrentPlayer()
	}
	if ws := g.getWinnerSymbol(); ws != noWinner {
		winner := g.getPlayerWithSymbol(ws)
		winner.incrementScore()
	}
	if g.shouldResetBoard() {
		g.resetBoard()
	}
	g.dispatchGame()
	g.dispatchNextTurn()
}

func rejectConnection(conn net.Conn) {
	fmt.Printf("rejecting client connected from %v\n", conn.RemoteAddr().String())
	conn.Write([]byte("Game is full, try again later!" + "\n"))
	conn.Close()
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

		if game.isFullPlaces() {
			rejectConnection(conn)
			continue
		}

		fmt.Printf("client connected from %v\n", conn.RemoteAddr().String())
		player := Player{
			Index:      len(game.Players),
			Connection: conn,
			Symbol:     []string{xSymbol, oSymbol}[len(game.Players)],
			Score:      0,
		}
		game.Players = append(game.Players, player)
		game.CurrentPlayer = player

		if game.isFullPlaces() {
			game.dispatchGame()
			game.dispatchNextTurn()
		}

		go handleGameConnection(&game, player)
	}
}
