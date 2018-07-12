package profiling

import (
	"strings"
	"testing"
)

var split []string

func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		split = strings.Split("one,two,three,four,five,six,seven", ",")
	}
}
