package openrt

import "unsafe"

func Memclr(p []byte) {
	if len(p) == 0 {
		return
	}
	ptr := unsafe.Pointer(&p[0])
	memclrNoHeapPointers(ptr, uintptr(len(p)))
}

//go:noescape
//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)
