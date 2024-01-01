package simdbyte

import "golang.org/x/sys/cpu"

var avxAllow bool

func init() {
	avxAllow = cpu.X86.HasAVX2 && cpu.X86.HasAVX512VPCLMULQDQ
}
