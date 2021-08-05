package cmhp_crypto_test

import (
	"fmt"
	"testing"

	"github.com/maldan/go-cmhp/cmhp_crypto"
)

func TestOne(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(cmhp_crypto.UID(10))
	}
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmhp_crypto.UID(10)
	}
}
