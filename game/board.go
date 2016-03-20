package game

import (
	"bytes"
	"errors"
	"fmt"

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

func (b *boardState) floodFill(root *bit.Vector, bounds *bit.Vector) *bit.Vector {
	for {
		next := root.Copy()
		next.Or(root.Copy().Lsh(1).AndNot(b.g.r))
		next.Or(root.Copy().Rsh(1).AndNot(b.g.l))
		next.Or(root.Copy().Lsh(uint(b.g.size)))
		next.Or(root.Copy().Rsh(uint(b.g.size)))
		next.AndNot(bounds)
		if next.Equal(root) {
			break
		}
		root = next
	}
	return root
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

func (b *boardState) String() string {
	out := &bytes.Buffer{}
	for r := 0; r < b.g.size; r++ {
		fmt.Fprintf(out, "% 2d", r)
		for c := 0; c < b.g.size; c++ {
			c, ok := b.at(c, r)
			switch {
			case !ok:
				fmt.Fprintf(out, " +")
			case c == White:
				fmt.Fprintf(out, " O")
			case c == Black:
				fmt.Fprintf(out, " X")
			}
		}
		fmt.Fprintf(out, "\n")
	}
	return out.String()
}
