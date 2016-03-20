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
	i, mask := v.pos(bit)
	return v.data[i]&mask != 0
}

// Set sets the specified bit
func (v *Vector) Set(bit int) {
	i, mask := v.pos(bit)
	v.data[i] |= mask
}

// Clear clears the specified bit
func (v *Vector) Clear(bit int) {
	i, mask := v.pos(bit)
	v.data[i] &= ^mask
}

// Or logically-or's the rhs into this vector
func (v *Vector) Or(rhs *Vector) {
	if v.Len() != rhs.Len() {
		panic("Or(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] |= w
	}
}

// And logically-and's the rhs into this vector
func (v *Vector) And(rhs *Vector) {
	if v.Len() != rhs.Len() {
		panic("And(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] &= w
	}
}

// Xor logically-xor's the rhs into this vector
func (v *Vector) Xor(rhs *Vector) {
	if v.Len() != rhs.Len() {
		panic("Xor(): len mismatch")
	}
	for i, w := range rhs.data {
		v.data[i] ^= w
	}
}

// Not negates the value of every bit in v
func (v *Vector) Not() {
	for i, w := range v.data {
		v.data[i] = ^w
	}
}

// Lsh shifts the input vector `bits` positions left (towards lower
// indexes)
func (v *Vector) Lsh(bits uint) {
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
func (v *Vector) Rsh(bits uint) {
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
