package openrt

import "unsafe"

// Memequal copies n bytes from "from" to "to".
func Memequal(a, b unsafe.Pointer, size uintptr) bool {
	return memequal(a, b, size)
}

//go:noescape
//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool
