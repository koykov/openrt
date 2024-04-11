package openrt

import "unsafe"

// Memmove copies n bytes from "from" to "to".
func Memmove(to, from unsafe.Pointer, n uintptr) {
	memmove(to, from, n)
}

//go:noescape
//go:linkname memmove runtime.memmove
func memmove(to, from unsafe.Pointer, n uintptr)
