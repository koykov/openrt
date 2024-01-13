package simdbyte

import (
	"bytes"
	"fmt"
	"testing"
)

type tc struct {
	a string
	b string
	i int
}

var indexTC = []tc{
	{"", "", 0},
	{"", "a", -1},
	{"", "foo", -1},
	{"fo", "foo", -1},
	{"foo", "baz", -1},
	{"foo", "foo", 0},
	{"oofofoofooo", "f", 2},
	{"oofofoofooo", "foo", 4},
	{"barfoobarfoo", "foo", 3},
	{"foo", "", 0},
	{"foo", "o", 1},
	{"abcABCabc", "A", 3},
	// cases with one byte strings - test IndexByte and special case in Index()
	{"", "a", -1},
	{"x", "a", -1},
	{"x", "x", 0},
	{"abc", "a", 0},
	{"abc", "b", 1},
	{"abc", "c", 2},
	{"abc", "x", -1},
	{"barfoobarfooyyyzzzyyyzzzyyyzzzyyyxxxzzzyyy", "x", 33},
	{"fofofofooofoboo", "oo", 7},
	{"fofofofofofoboo", "ob", 11},
	{"fofofofofofoboo", "boo", 12},
	{"fofofofofofoboo", "oboo", 11},
	{"fofofofofoooboo", "fooo", 8},
	{"fofofofofofoboo", "foboo", 10},
	{"fofofofofofoboo", "fofob", 8},
	{"fofofofofofofoffofoobarfoo", "foffof", 12},
	{"fofofofofoofofoffofoobarfoo", "foffof", 13},
	{"fofofofofofofoffofoobarfoo", "foffofo", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofo", 13},
	{"fofofofofoofofoffofoobarfoo", "foffofoo", 13},
	{"fofofofofofofoffofoobarfoo", "foffofoo", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofoob", 13},
	{"fofofofofofofoffofoobarfoo", "foffofoob", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofooba", 13},
	{"fofofofofofofoffofoobarfoo", "foffofooba", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofoobar", 13},
	{"fofofofofofofoffofoobarfoo", "foffofoobar", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofoobarf", 13},
	{"fofofofofofofoffofoobarfoo", "foffofoobarf", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofoobarfo", 13},
	{"fofofofofofofoffofoobarfoo", "foffofoobarfo", 12},
	{"fofofofofoofofoffofoobarfoo", "foffofoobarfoo", 13},
	{"fofofofofofofoffofoobarfoo", "foffofoobarfoo", 12},
	{"fofofofofoofofoffofoobarfoo", "ofoffofoobarfoo", 12},
	{"fofofofofofofoffofoobarfoo", "ofoffofoobarfoo", 11},
	{"fofofofofoofofoffofoobarfoo", "fofoffofoobarfoo", 11},
	{"fofofofofofofoffofoobarfoo", "fofoffofoobarfoo", 10},
	{"fofofofofoofofoffofoobarfoo", "foobars", -1},
	{"foofyfoobarfoobar", "y", 4},
	{"oooooooooooooooooooooo", "r", -1},
	{"oxoxoxoxoxoxoxoxoxoxoxoy", "oy", 22},
	{"oxoxoxoxoxoxoxoxoxoxoxox", "oy", -1},
	// test fallback to Rabin-Karp.
	{"000000000000000000000000000000000000000000000000000000000000000000000001", "0000000000000000000000000000000000000000000000000000000000000000001", 5},
}

func TestIndexSSE(t *testing.T) {
	for _, tc_ := range indexTC {
		t.Run(fmt.Sprintf("%s_%s", tc_.a, tc_.b), func(t *testing.T) {
			i := IndexSSE([]byte(tc_.a), []byte(tc_.b))
			if i != tc_.i {
				t.FailNow()
			}
		})
	}
}

func TestIndexByteSSE(t *testing.T) {
	for _, tc_ := range indexTC {
		if len(tc_.b) > 1 || len(tc_.b) == 0 {
			continue
		}
		t.Run(fmt.Sprintf("%s__%s", tc_.a, tc_.b), func(t *testing.T) {
			i := IndexByteSSE([]byte(tc_.a), tc_.b[0])
			if i != tc_.i {
				t.FailNow()
			}
		})
	}
}

func BenchmarkIndexByteSSE(b *testing.B) {
	for _, tc_ := range indexTC {
		if len(tc_.b) > 1 || len(tc_.b) == 0 {
			continue
		}
		b.Run(fmt.Sprintf("sse: %s__%s", tc_.a, tc_.b), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				x := IndexByteSSE([]byte(tc_.a), tc_.b[0])
				if x != tc_.i {
					b.FailNow()
				}
			}
		})
		b.Run(fmt.Sprintf("vanilla: %s__%s", tc_.a, tc_.b), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				x := bytes.IndexByte([]byte(tc_.a), tc_.b[0])
				if x != tc_.i {
					b.FailNow()
				}
			}
		})
	}
}
