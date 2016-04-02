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

	// ErrKo is returned if a move is illegal because it repeats
	// the previous position
	ErrKo = errors.New("move results in an illegal ko capture")

	// ErrGameOver is returned if a move is played on a game that
	// has completed
	ErrGameOver = errors.New("game is over")
)

// boardState represents a specific game position. It is immutable
// once created
type boardState struct {
	g                              *Game
	prev                           *boardState
	white, black                   *bit.Vector
	blackPrisoners, whitePrisoners int
	toPlay                         Color
	passes                         int
}

func (b *boardState) move(x, y int) (*boardState, error) {
	if x < 0 || x >= b.g.Size || y < 0 || y >= b.g.Size {
		return nil, ErrOutOfBounds
	}
	idx := y*b.g.Size + x
	if b.white.At(idx) || b.black.At(idx) {
		return nil, ErrOccupied
	}
	out := *b
	out.prev = b
	var me, them **bit.Vector
	var prisoners *int
	if b.toPlay == White {
		me, them = &out.white, &out.black
		prisoners = &out.whitePrisoners
	} else {
		them, me = &out.white, &out.black
		prisoners = &out.blackPrisoners
	}
	*me = (*me).Copy().Set(idx)

	if x > 0 {
		if c := out.deadGroupAt(idx-1, *them, *me); c != nil {
			*them = (*them).Copy().AndNot(c)
			*prisoners += c.Popcount()
		}
	}
	if x < b.g.Size-1 {
		if c := out.deadGroupAt(idx+1, *them, *me); c != nil {
			*them = (*them).Copy().AndNot(c)
			*prisoners += c.Popcount()
		}
	}
	if y > 0 {
		if c := out.deadGroupAt(idx-b.g.Size, *them, *me); c != nil {
			*them = (*them).Copy().AndNot(c)
			*prisoners += c.Popcount()
		}
	}
	if y < b.g.Size-1 {
		if c := out.deadGroupAt(idx+b.g.Size, *them, *me); c != nil {
			*them = (*them).Copy().AndNot(c)
			*prisoners += c.Popcount()
		}
	}
	if c := out.deadGroupAt(idx, *me, *them); c != nil {
		return nil, ErrSelfCapture
	}

	if b.prev != nil && b.prev.white.Equal(out.white) &&
		b.prev.black.Equal(out.black) {
		return nil, ErrKo
	}

	out.toPlay = !out.toPlay
	out.passes = 0
	return &out, nil
}

func (b *boardState) pass() *boardState {
	next := *b
	next.prev = b
	next.toPlay = !next.toPlay
	next.passes++
	return &next
}

func (b *boardState) gameOver() bool {
	return b.passes >= 2
}

func (b *boardState) deadGroupAt(idx int, me *bit.Vector, them *bit.Vector) *bit.Vector {
	group := b.floodFill(bit.NewVector(b.white.Len()).Set(idx),
		me.Copy().Not())
	if b.grow(group).AndNot(group).AndNot(them).Popcount() == 0 {
		return group
	}
	return nil
}

func (b *boardState) grow(root *bit.Vector) *bit.Vector {
	next := root.Copy()
	next.Or(root.Copy().Lsh(1).AndNot(b.g.r))
	next.Or(root.Copy().Rsh(1).AndNot(b.g.l))
	next.Or(root.Copy().Lsh(uint(b.g.Size)))
	next.Or(root.Copy().Rsh(uint(b.g.Size)))
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
	bit := y*b.g.Size + x
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
	for r := 0; r < b.g.Size; r++ {
		fmt.Fprintf(out, "% 2d", r)
		for c := 0; c < b.g.Size; c++ {
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
