package openrt

import "unsafe"

// Noescape hides a pointer from escape analysis.
func Noescape(ptr unsafe.Pointer) unsafe.Pointer {
	return noescape(ptr)
}

//go:noescape
//go:linkname noescape runtime.noescape
func noescape(ptr unsafe.Pointer) unsafe.Pointer

var _ = Noescape
