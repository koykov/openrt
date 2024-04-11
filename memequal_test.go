package openrt

import (
	"math"
	"testing"
	"unsafe"
)

const meqMB = 1024 * 1024

var meqA, meqB []byte

func init() {
	meqA, meqB = make([]byte, meqMB), make([]byte, meqMB)
	for i := 0; i < meqMB; i++ {
		meqA[i], meqB[i] = byte(i%math.MaxUint8), byte(i%math.MaxUint8)
	}
}

func TestMemequal(t *testing.T) {
	if !Memequal(unsafe.Pointer(&meqA[0]), unsafe.Pointer(&meqB[0]), meqMB) {
		t.FailNow()
	}
}

func BenchmarkMemequal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Memequal(unsafe.Pointer(&meqA[0]), unsafe.Pointer(&meqB[0]), meqMB)
	}
}
