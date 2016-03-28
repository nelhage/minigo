package sgf

import (
	"reflect"
	"strings"
	"testing"
)

const basic = `
(;FF[4]GM[1]SZ[19]
 GN[Copyright goproblems.com]
 PB[Black]
 HA[0]
 PW[White]
 KM[5.5]
 DT[1999-07-21]
 TM[1800]
 RU[Japanese]
 ;AW[bb][cb][cc][cd][de][df][cg][ch][dh][ai][bi][ci]
 AB[ba][ab][ac][bc][bd][be][cf][bg][bh]
 C[Black to play and live.]
 (;B[af];W[ah]
 (;B[ce];W[ag]C[only one eye this way])
 (;B[ag];W[ce]))
 (;B[ah];W[af]
 (;B[ae];W[bf];B[ag];W[bf]
 (;B[af];W[ce]C[oops! you can't take this stone])
 (;B[ce];W[af];B[bg]C[RIGHT black plays under the stones and lives]))
 (;B[bf];W[ae]))
 (;B[ae];W[ag]))
`

func TestBasic(t *testing.T) {
	c, e := ParseSGF(strings.NewReader(basic))
	if e != nil {
		t.Fatal("parse error", e)
	}
	if len(c.Trees) != 1 {
		t.Error("wrong number of trees")
	}
	g := c.Trees[0]
	if len(g.Principal.Nodes) != 2 {
		t.Error("wrong number of nodes", g.Principal.Nodes)
	}
	var names []string
	for _, p := range g.Principal.Nodes[0].Props {
		names = append(names, p.Prop)
	}
	expect := []string{"FF", "GM", "SZ", "GN", "PB", "HA", "PW", "KM", "DT", "TM", "RU"}
	if !reflect.DeepEqual(names, expect) {
		t.Errorf("expect %#v got %#v",
			expect, names)
	}
}
