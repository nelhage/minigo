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

	// ErrSelfCapture is returned if a requested move would result
	// in capture of the placed stone
	ErrSelfCapture = errors.New("requested move results in self-capture")
)

// boardState represents a specific game position. It is immutable
// once created
type boardState struct {
	g      *Game
	prev   *boardState
	white  *bit.Vector
	black  *bit.Vector
	toPlay Color
}

func (b *boardState) move(x, y int) (*boardState, error) {
	if x < 0 || x >= b.g.size || y < 0 || y >= b.g.size {
		return nil, ErrOutOfBounds
	}
	idx := y*b.g.size + x
	if b.white.At(idx) || b.black.At(idx) {
		return nil, ErrOccupied
	}
	out := *b
	out.prev = b
	var me, them **bit.Vector
	if b.toPlay == White {
		me, them = &out.white, &out.black
	} else {
		them, me = &out.white, &out.black
	}
	*me = (*me).Copy().Set(idx)
	group := out.floodFill(bit.NewVector(b.white.Len()).Set(idx),
		(*me).Copy().Not())
	if b.grow(group).AndNot(group).AndNot(*them).Popcount() == 0 {
		return nil, ErrSelfCapture
	}

	out.toPlay = !out.toPlay
	return &out, nil
}

func (b *boardState) grow(root *bit.Vector) *bit.Vector {
	next := root.Copy()
	next.Or(root.Copy().Lsh(1).AndNot(b.g.r))
	next.Or(root.Copy().Rsh(1).AndNot(b.g.l))
	next.Or(root.Copy().Lsh(uint(b.g.size)))
	next.Or(root.Copy().Rsh(uint(b.g.size)))
	return next
}

func (b *boardState) floodFill(root *bit.Vector, bounds *bit.Vector) *bit.Vector {
	for {
		next := b.grow(root)
		next.AndNot(bounds)
		if next.Equal(root) {
			break
		}
		root = next
	}
	return root
}

func (b *boardState) at(x, y int) (Color, bool) {
	bit := y*b.g.size + x
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
