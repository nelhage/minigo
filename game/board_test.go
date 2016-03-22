package game

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"nelhage.com/minigo/bit"
)

var fixtureRE = regexp.MustCompile(`\A\s*((?:\d+\s*(?:[OX+*]\s*)+\n)+)`)

func board(g *Game, in string) *boardState {
	white := bit.NewVector(g.size * g.size)
	black := bit.NewVector(g.size * g.size)

	m := fixtureRE.FindStringSubmatch(in)
	if m == nil {
		panic(fmt.Sprintf("bad fixture:\n%s", in))
	}
	lines := strings.Split(strings.TrimRight(m[1], "\n"), "\n")
	if len(lines) != g.size {
		panic(fmt.Sprintf("bad fixture (%d rows):\n%#v", len(lines), lines))
	}
	for i, l := range lines {
		bits := strings.Split(l, " ")
		if len(bits) != g.size+1 {
			panic(fmt.Sprintf("bad fixture:\n%s", in))
		}
		for j, c := range bits[1:] {
			switch c {
			case "*", "+":
			case "O":
				white.Set(g.size*i + j)
			case "X":
				black.Set(g.size*i + j)
			}
		}
	}
	return &boardState{
		g:     g,
		white: white,
		black: black,
	}
}

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
	cases := []struct{ in, out string }{
		{`
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * + + + * + +
3 + + + + + + + + +
4 + + + + O + + + +
5 + + + + + + + + +
6 + + * + + + * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`, `
0 O O O O O O O O O
1 O O O O O O O O O
2 O O O O O O O O O
3 O O O O O O O O O
4 O O O O O O O O O
5 O O O O O O O O O
6 O O O O O O O O O
7 O O O O O O O O O
8 O O O O O O O O O
  0 1 2 3 4 5 6 7 8
`,
		},
		{`
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * X X X * + +
3 + + X + + + X + +
4 + + X + O + X + +
5 + + X + + + X + +
6 + + * X X X * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`, `
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * X X X * + +
3 + + X O O O X + +
4 + + X O O O X + +
5 + + X O O O X + +
6 + + * X X X * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
		},
		{`
0 + + + X O X + + +
1 + + + + X + + + +
2 + + * + + + * + +
3 X + + + + + + + X
4 O X + + + + + X O
5 X + + + + + + + X
6 + + * X + + * + +
7 + + X + X + + + +
8 + X + O + X + + +
  0 1 2 3 4 5 6 7 8
`, `
0 + + + X O X + + +
1 + + + + X + + + +
2 + + * + + + * + +
3 X + + + + + + + X
4 O X + + + + + X O
5 X + + + + + + + X
6 + + * X + + * + +
7 + + X O X + + + +
8 + X O O O X + + +
  0 1 2 3 4 5 6 7 8
`,
		},
	}
	for _, tc := range cases {
		g := New(9)
		in := board(g, tc.in)
		out := board(g, tc.out)

		fill := in.floodFill(in.white, in.black)
		if !fill.Equal(out.white) {
			t.Logf("fill=%#v want=%#v", fill, out.white)
			t.Logf("count(fill)=%d count(want)=%d", fill.Popcount(), out.white.Popcount())
			t.Errorf("wrong fill in=\n%s\nwant=\n%s\nout=\n%s",
				in, out,
				&boardState{white: fill, black: out.black, g: g})
		}
	}
}

func TestSelfKill(t *testing.T) {
	g := New(9)
	g.board = board(g, `
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * + + + * + +
3 + + + + X + + + +
4 + + + X + X + + +
5 + + + + X + + + +
6 + + * + + + * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`)
	g.board.toPlay = White
	if err := g.Move(4, 4); err != ErrSelfCapture {
		t.Fatal("game allowed self-capture")
	}
}

func TestCapture(t *testing.T) {
	cases := []struct {
		in   string
		who  Color
		x, y int
		out  string
	}{
		{
			`
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * + + + * + +
3 + + + + + + + + +
4 + + + X O X + + +
5 + + + + X + + + +
6 + + * + + + * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
			Black, 4, 3,
			`
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * + + + * + +
3 + + + + X + + + +
4 + + + X + X + + +
5 + + + + X + + + +
6 + + * + + + * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
		},
		{
			`
0 + + + + X X X + +
1 + + + X O O O O X
2 + + * X O X X X +
3 + + + X O X + + +
4 + + + X O X + + +
5 + + + + X + + + +
6 + + * + + + * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
			Black, 7, 0,
			`
0 + + + + X X X X +
1 + + + X + + + + X
2 + + * X + X X X +
3 + + + X + X + + +
4 + + + X + X + + +
5 + + + + X + + + +
6 + + * + + + * + +
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
		},
		{
			`
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * + + + * + +
3 + + + + + + + + O
4 + + + + + + + O X
5 + + + + + + + O X
6 + + * + + + * O X
7 + + + + + + + + +
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
			White, 8, 7,
			`
0 + + + + + + + + +
1 + + + + + + + + +
2 + + * + + + * + +
3 + + + + + + + + O
4 + + + + + + + O +
5 + + + + + + + O +
6 + + * + + + * O +
7 + + + + + + + + O
8 + + + + + + + + +
  0 1 2 3 4 5 6 7 8
`,
		},
	}
	for i, tc := range cases {
		g := New(9)
		g.board = board(g, tc.in)
		g.board.toPlay = tc.who
		if err := g.Move(tc.x, tc.y); err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		expect := board(g, tc.out)
		if !expect.white.Equal(g.board.white) ||
			!expect.black.Equal(g.board.black) {
			t.Errorf("%d: want:\n%s\ngot:\n%s", i, expect, g.board)
		}
	}
}
