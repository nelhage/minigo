package bit

// Vector implements a fixed-width bit vector
type Vector struct {
	data []uint64
	bits int
}

// NewVector returns a new zero-initialized bit vector of length
// `bits` bits.
func NewVector(bits int) *Vector {
	return &Vector{
		data: make([]uint64, (bits+63)/64),
		bits: bits,
	}
}

func (v *Vector) normalize() {
	trailing := uint(v.bits % 64)
	if trailing != 0 {
		mask := (uint64(1) << trailing) - 1
		v.data[len(v.data)-1] &= mask
	}
}

func (v *Vector) checkNormalized() {
	trailing := uint(v.bits % 64)
	if trailing != 0 {
		mask := (uint64(1) << trailing) - 1
		if v.data[len(v.data)-1]&mask != v.data[len(v.data)-1] {
			panic("checkNormalized")
		}
	}
}

// Copy returns a new vector with identical length and bit pattern,
// which does not share any storage with the original vector
func (v *Vector) Copy() *Vector {
	newdata := make([]uint64, len(v.data))
	copy(newdata, v.data)
	return &Vector{
		data: newdata,
		bits: v.bits,
	}
}

// Bytes returns the bit vector as an array of bytes in LSB order. The
// returned array is a copy of the underlying data.
func (v *Vector) Bytes() []byte {
	out := make([]byte, v.bits/8)
	for i := 0; i < len(out); i++ {
		out[i] = byte(v.data[i/4] >> (uint(i) % 4 * 8))
	}
	return out
}

// Len returns the length of the vector in bits
func (v *Vector) Len() int {
	return v.bits
}

func (v *Vector) pos(bit int) (i int, mask uint64) {
	return bit / 64, uint64(1) << (uint(bit) % 64)
}

// At returns the value of the bit at position `bit`
func (v *Vector) At(bit int) bool {
	if bit >= v.bits {
		panic("at: out of range")
	}
	i, mask := v.pos(bit)
	return v.data[i]&mask != 0
}

// Set sets the specified bit and returns the input vector
func (v *Vector) Set(bit int) *Vector {
	if bit >= v.bits {
		panic("set: out of range")
	}
	i, mask := v.pos(bit)
	v.data[i] |= mask
	return v
}

// Clear clears the specified bit and returns the input vector
func (v *Vector) Clear(bit int) *Vector {
	if bit >= v.bits {
		panic("clear: out of range")
	}
	i, mask := v.pos(bit)
	v.data[i] &= ^mask
	return v
}

// Or logically-or's the rhs into this vector and returns the input
// vector
func (v *Vector) Or(rhs *Vector) *Vector {
	if v.Len() != rhs.Len() {
		panic("Or(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] |= w
	}
	v.normalize()
	return v
}

// And logically-and's the rhs into this vector
func (v *Vector) And(rhs *Vector) *Vector {
	if v.Len() != rhs.Len() {
		panic("And(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] &= w
	}
	return v
}

// AndNot is equivalent to v.And(rhs.Copy().Not()), but saves a copy.
func (v *Vector) AndNot(rhs *Vector) *Vector {
	if v.Len() != rhs.Len() {
		panic("AndNot(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] &= ^w
	}
	return v
}

// Xor logically-xor's the rhs into this vector
func (v *Vector) Xor(rhs *Vector) *Vector {
	if v.Len() != rhs.Len() {
		panic("Xor(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] ^= w
	}
	v.normalize()
	return v
}

// Not negates the value of every bit in v
func (v *Vector) Not() *Vector {
	for i, w := range v.data {
		v.data[i] = ^w
	}
	v.normalize()
	return v
}

// Lsh shifts the input vector `bits` positions left (towards lower
// indexes)
func (v *Vector) Lsh(bits uint) *Vector {
	full := bits / 64
	partial := bits % 64
	if full != 0 {
		v.fullShiftLeft(full)
	}
	var accum uint64
	for i := len(v.data) - 1; i >= 0; i-- {
		out := v.data[i] << (64 - partial)
		v.data[i] = v.data[i]>>partial | accum
		accum = out
	}
	return v
}

// Shift the words in `data` `off` positions towards lower indexes
func (v *Vector) fullShiftLeft(off uint) {
	for i := range v.data {
		if i+int(off) < len(v.data) {
			v.data[i] = v.data[uint(i)+off]
		} else {
			v.data[i] = 0
		}
	}
}

// Rsh right-shifts vector by `bits` bits (towards higher indexes)
func (v *Vector) Rsh(bits uint) *Vector {
	full := bits / 64
	partial := bits % 64
	if full != 0 {
		v.fullShiftRight(full)
	}
	var accum uint64
	for i := 0; i < len(v.data); i++ {
		out := v.data[i] >> (64 - partial)
		v.data[i] = v.data[i]<<partial | accum
		accum = out
	}
	v.normalize()
	return v
}

// Shift the words in `data` `off` positions towards higher indexes
func (v *Vector) fullShiftRight(off uint) {
	for i := len(v.data) - 1; i >= 0; i-- {
		if i-int(off) >= 0 {
			v.data[i] = v.data[uint(i)-off]
		} else {
			v.data[i] = 0
		}
	}
}

// Equal returns true iff the lhs and rhs have identical bit patterns
func (v *Vector) Equal(rhs *Vector) bool {
	v.checkNormalized()
	rhs.checkNormalized()

	for i, w := range v.data {
		if rhs.data[i] != w {
			return false
		}
	}
	return true
}

// Popcount returns the number of bits set in `v`
func (v *Vector) Popcount() int {
	p := 0
	for _, w := range v.data {
		p += popcount(w)
	}
	return p
}

func popcount(x uint64) (n int) {
	// bit population count, see
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return int(x >> 56)
}
