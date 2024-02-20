package openrt

import "unsafe"

func Memclr(p []byte) {
	if len(p) == 0 {
		return
	}
	MemclrUnsafe(unsafe.Pointer(&p[0]), len(p))
}

func MemclrUnsafe(ptr unsafe.Pointer, len_ int) {
	memclrNoHeapPointers(ptr, uintptr(len_))
}

//go:noescape
//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)
