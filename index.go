package simdbyte

import "golang.org/x/sys/cpu"

func init() {
	_ = cpu.X86.HasSSE42
}
