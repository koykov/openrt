package openrt

import "unsafe"

// Memequal checks equality of a and b. Both arrays must have same length.
func Memequal(a, b unsafe.Pointer, size uintptr) bool {
	return memequal(a, b, size)
}

//go:noescape
//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool
