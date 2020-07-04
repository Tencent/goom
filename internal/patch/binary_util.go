package patch

// A ByteOrder specifies how to convert byte sequences into
// 16-, 32-, or 64-bit unsigned integers.
type ByteOrder interface {
	Uint16([]byte) uint16
	Uint32([]byte) uint32
	Uint64([]byte) uint64
	PutUint16([]byte, uint16)
	PutUint32([]byte, uint32)
	PutUint64([]byte, uint64)
	String() string
}

// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian

type littleEndian struct{}

func (littleEndian) Int16(b []byte) int16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return int16(b[0]) | int16(b[1])<<8
}

func (littleEndian) PutInt16(b []byte, v int16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func (littleEndian) Int32(b []byte) int32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

func (littleEndian) PutInt32(b []byte, v int32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func (littleEndian) Int64(b []byte) int64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 |
		int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
}

func (littleEndian) PutInt64(b []byte, v int64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func (littleEndian) String() string { return "LittleEndian" }

func (littleEndian) GoString() string { return "binary.LittleEndian" }

type bigEndian struct{}

func (bigEndian) Int16(b []byte) int16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return int16(b[1]) | int16(b[0])<<8
}

func (bigEndian) PutInt16(b []byte, v int16) {
	_ = b[1] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func (bigEndian) Int32(b []byte) int32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

func (bigEndian) PutInt32(b []byte, v int32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func (bigEndian) Int64(b []byte) int64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return int64(b[7]) | int64(b[6])<<8 | int64(b[5])<<16 | int64(b[4])<<24 |
		int64(b[3])<<32 | int64(b[2])<<40 | int64(b[1])<<48 | int64(b[0])<<56
}

func (bigEndian) PutInt64(b []byte, v int64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

func (bigEndian) String() string { return "BigEndian" }

func (bigEndian) GoString() string { return "binary.BigEndian" }
