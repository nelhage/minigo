package game

import "nelhage.com/minigo/bit"

// Color represents one of the two players in Go
type Color bool

var (
	// White is the white player
	White Color = true
	// Black is the Black player
	Black Color = false
)

// Game represents a game of Go
type Game struct {
	Size  int
	board *boardState

	l, r, t, b *bit.Vector
	z          *bit.Vector
}

// New returns a new game of board size `size` on a side
func New(size int) *Game {
	g := &Game{Size: size}
	g.board = &boardState{
		g:      g,
		white:  bit.NewVector(size * size),
		black:  bit.NewVector(size * size),
		toPlay: Black,
	}
	g.precompute()
	return g
}

func (g *Game) precompute() {
	g.l = bit.NewVector(g.Size * g.Size)
	g.r = bit.NewVector(g.Size * g.Size)
	g.t = bit.NewVector(g.Size * g.Size)
	g.b = bit.NewVector(g.Size * g.Size)
	g.z = bit.NewVector(g.Size * g.Size)
	for i := 0; i < g.Size; i++ {
		g.l.Set(i * g.Size)
		g.r.Set((i+1)*g.Size - 1)
		g.t.Set(i)
		g.b.Set(g.Size*(g.Size-1) + i)
	}
}

// ToPlay returns the player whose turn it is
func (g *Game) ToPlay() Color {
	return g.board.toPlay
}

// GameOver returns true if the game is over
func (g *Game) GameOver() bool {
	return g.board.gameOver()
}

// Move plays a stone at position (x,y)
func (g *Game) Move(x, y int) error {
	if g.board.gameOver() {
		return ErrGameOver
	}
	b, err := g.board.move(x, y)
	if err != nil {
		return err
	}
	g.board = b
	return nil
}

// Pass causes the current player to pass
func (g *Game) Pass() error {
	if g.board.gameOver() {
		return ErrGameOver
	}
	g.board = g.board.pass()
	return nil
}

// At returns a boolean indicating whether a given intersection is
// populated, and the color of the stone at that intersection if there
// is one
func (g *Game) At(x, y int) (Color, bool) {
	return g.board.at(x, y)
}
