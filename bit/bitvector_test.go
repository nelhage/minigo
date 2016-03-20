package bit

import (
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		len, outlen int
	}{
		{0, 0},
		{1, 1},
		{64, 1},
		{65, 2},
		{127, 2},
	}
	for _, tc := range cases {
		v := NewVector(tc.len)
		if v.bits != tc.len {
			t.Errorf("len(v)=%d != %d", v.bits, tc.len)
		}
		if len(v.data) != tc.outlen {
			t.Errorf("New(%d): len(data)=%d != %d",
				tc.len, len(v.data), tc.outlen)
		}
		for i, w := range v.data {
			if w != 0 {
				t.Errorf("New(%d): [%d]=%d!=0", tc.len, i, w)
			}
		}
	}
}

func TestAt(t *testing.T) {
	v := NewVector(256)
	for i := 0; i < v.Len(); i++ {
		if v.At(i) {
			t.Error("bit ", i, " set")
		}
		v.Set(i)
		if !v.At(i) {
			t.Error("Set(", i, ") did not set")
		}
		v.Clear(i)
		if v.At(i) {
			t.Error("Clear(", i, ") did not clear")
		}
	}
}

func TestLsh(t *testing.T) {
	for b := 0; b < 256; b++ {
		b := 128
		for s := uint(0); s < 256; s++ {
			v := NewVector(256)
			v.Set(b)
			v.Lsh(s)
			for i := 0; i < 256; i++ {
				if (i == b-int(s)) != v.At(i) {
					t.Errorf("bit(%d)<<%d [%d]=%v",
						b, s, i, v.At(i),
					)
				}
			}
		}
	}
}

func TestRsh(t *testing.T) {
	for b := 0; b < 256; b++ {
		b := 128
		for s := uint(0); s < 256; s++ {
			v := NewVector(256)
			v.Set(b)
			v.Rsh(s)
			for i := 0; i < 256; i++ {
				if (i == b+int(s)) != v.At(i) {
					t.Errorf("bit(%d)<<%d [%d]=%v",
						b, s, i, v.At(i),
					)
				}
			}
		}
	}
}

func TestUnaligned(t *testing.T) {
	a := NewVector(100)
	b := NewVector(100)
	a.Set(99)
	a.Rsh(2)
	if !a.Equal(b) {
		t.Errorf("rsh into undefined bits lived: %x", a.data[len(a.data)-1])
	}
}
