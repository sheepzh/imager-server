package corpus

import (
	"testing"
)

func TestPageQuery(b *testing.T) {
}

func BenchMarkListAuthors(b *testing.B) {
	for n := 0; n < b.N; n++ {
		__ListAuthors(1, 2000)
	}
}
