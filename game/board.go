package game

import (
	"errors"

	"nelhage.com/minigo/bit"
)

var (
	// ErrOutOfBounds is returned if a requested position is
	// outside the bounds of the board
	ErrOutOfBounds = errors.New("position out of bounds")

	// ErrOccupied is returned if a requested move is on top of an
	// existing stone.
	ErrOccupied = errors.New("requested position is already occupied")
)

// boardState represents a specific game position. It is immutable
// once created
type boardState struct {
	g      *Game
	white  *bit.Vector
	black  *bit.Vector
	toPlay Color
}

func (b *boardState) move(x, y int) (*boardState, error) {
	if x < 0 || x >= b.g.size || y < 0 || y >= b.g.size {
		return nil, ErrOutOfBounds
	}
	bit := x*b.g.size + y
	if b.white.At(bit) || b.black.At(bit) {
		return nil, ErrOccupied
	}
	out := *b
	if b.toPlay == White {
		out.white = out.white.Copy().Set(bit)
	} else {
		out.black = out.black.Copy().Set(bit)
	}
	out.toPlay = !out.toPlay
	return &out, nil
}

func (b *boardState) at(x, y int) (Color, bool) {
	bit := x*b.g.size + y
	if b.white.At(bit) {
		return White, true
	}
	if b.black.At(bit) {
		return Black, true
	}
	return Black, false
}
