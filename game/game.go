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
	size  int
	board *boardState
}

// New returns a new game of board size `size` on a side
func New(size int) *Game {
	g := &Game{size: size}
	g.board = &boardState{
		g:      g,
		white:  bit.NewVector(size * size),
		black:  bit.NewVector(size * size),
		toPlay: Black,
	}
	return g
}

// ToPlay returns the player whose turn it is
func (g *Game) ToPlay() Color {
	return g.board.toPlay
}

// Move plays a stone at position (x,y)
func (g *Game) Move(x, y int) error {
	b, err := g.board.move(x, y)
	if err != nil {
		return err
	}
	g.board = b
	return nil
}
