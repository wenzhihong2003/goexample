package profiling

import (
	"strings"
	"testing"
)

var upper string

func BenchmarkToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		upper = strings.ToUpper("Education is what remains after one has forgotten what one has learned in school.")
	}
}
