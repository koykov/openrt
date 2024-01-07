//go:build amd64

package simdbyte

import "bytes"

func IndexSSE(str, needle []byte) int {
	n, m := len(str), len(needle)
	if n == 0 || m == 0 || n < m {
		return -1
	}
	if n == m && bytes.Equal(str, needle) {
		return 0
	}
	if n < 16 || m < 12 {
		_ = str[n-1]
		for i := 0; i < n-len(needle); i++ {
			if bytes.Equal(str[i:i+m], needle) {
				return i
			}
		}
		return -1
	}
	return indexSSE(str, needle)
}

//go:noescape
func indexSSE(str, needle []byte) int

func IndexByteSSE(b []byte, c byte) int {
	if len(b) == 0 {
		return -1
	}
	return indexByteSSE(b, c)
}

//go:noescape
func indexByteSSE(b []byte, c byte) int
