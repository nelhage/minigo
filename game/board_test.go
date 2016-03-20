package game

import "testing"

func TestBasic(t *testing.T) {
	g := New(9)
	moves := []struct{ x, y int }{
		{4, 4}, {4, 2}, {4, 6}, {4, 7},
	}
	for _, m := range moves {
		err := g.Move(m.x, m.y)
		if err != nil {
			t.Fatalf("Move(%d,%d): %v", m.x, m.y, err)
		}
	}

	if err := g.Move(10, 10); err != ErrOutOfBounds {
		t.Fatalf("out of bounds move")
	}
	if g.ToPlay() != Black {
		t.Fatalf("bad move advanced player")
	}

	if err := g.Move(4, 7); err != ErrOccupied {
		t.Fatalf("occupied move")
	}

	ats := []struct {
		x, y int
		c    Color
		ok   bool
	}{
		{4, 4, Black, true},
		{4, 2, White, true},
		{4, 6, Black, true},
		{4, 7, White, true},
		{2, 4, Black, false},
		{2, 2, Black, false},
		{4, 8, Black, false},
	}
	for _, at := range ats {
		c, ok := g.At(at.x, at.y)
		if at.c != c || at.ok != ok {
			t.Errorf("At(%d,%d) = (%v,%v) != (%v, %v)",
				at.x, at.y, c, ok, at.c, at.ok)
		}
	}
}

func TestFloodFill(t *testing.T) {
	g := New(9)
	root := g.board.white.Copy().Set(4*9 + 4)
	flood := g.board.floodFill(root, g.board.white)
	if !flood.Equal(g.z.Copy().Not()) {
		b := &boardState{
			black: flood,
			white: g.board.white,
			g:     g,
		}
		t.Errorf("flood fill did not fill:\n%s", b.String())
	}

	bounds := g.board.white.Copy()
	for i := 2; i < 7; i++ {
		bounds.Set(2*g.size + i)
		bounds.Set(i*g.size + 2)
		bounds.Set(6*g.size + i)
		bounds.Set(i*g.size + 6)
	}

	flood = g.board.floodFill(root, bounds)
	b := &boardState{
		black: flood,
		white: bounds,
		g:     g,
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < g.size; j++ {
			check := []struct{ x, y int }{
				{i, j},
				{g.size - i, j},
				{j, i},
				{g.size - j, i},
			}
			for _, ch := range check {
				if _, ok := b.at(ch.x, ch.y); ok {
					t.Fatalf("At(%d,%d):\n%s",
						ch.x, ch.y, b)
				}
			}
		}
	}
}
