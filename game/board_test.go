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
