package openrt

import (
	"fmt"
	"strconv"
	"testing"
)

var memclrStages [][]byte

func init() {
	for i := 0; i < 20; i++ {
		l := 1 << i
		stg := make([]byte, l)
		for j := 0; j < len(stg); j++ {
			stg[j] = 1
		}
		memclrStages = append(memclrStages, stg)
	}
}

func TestMemclr(t *testing.T) {
	for i := 0; i < len(memclrStages); i++ {
		stg := memclrStages[i]
		t.Run(strconv.Itoa(len(stg)), func(t *testing.T) {
			Memclr(stg)
		})
	}
}

func BenchmarkMemclr(b *testing.B) {
	for i := 0; i < len(memclrStages); i++ {
		stg := memclrStages[i]
		b.Run(fmt.Sprintf("runtime/%d", len(stg)), func(b *testing.B) {
			b.ReportAllocs()
			for j := 0; j < b.N; j++ {
				Memclr(stg)
			}
		})
	}
}
