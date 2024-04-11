package openrt

import (
	"testing"
	"unsafe"
)

type mmove struct {
	a int
	b uint
	c float32
	d uint64
	e float64
	f int16
	g int64
}

func TestMemmove(t *testing.T) {
	a := mmove{
		a: 1,
		b: 2,
		c: 3,
		d: 4,
		e: 5,
		f: 6,
		g: 7,
	}
	var b mmove
	Memmove(unsafe.Pointer(&b), unsafe.Pointer(&a), unsafe.Sizeof(a))
	if b.d != 4 {
		t.FailNow()
	}
}

func BenchmarkMemmove(b *testing.B) {
	a := mmove{
		a: 1,
		b: 2,
		c: 3,
		d: 4,
		e: 5,
		f: 6,
		g: 7,
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var b_ mmove
		Memmove(unsafe.Pointer(&b_), unsafe.Pointer(&a), unsafe.Sizeof(a))
	}
}
