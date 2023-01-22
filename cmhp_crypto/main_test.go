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

func TestCrypt(t *testing.T) {
	for i := 0; i < 100; i++ {
		x, _ := cmhp_crypto.EncryptAes32([]byte(fmt.Sprintf("%v", i)), fmt.Sprintf("%v", i*10000))
		y, _ := cmhp_crypto.DecryptAes32(x, fmt.Sprintf("%v", i*10000))

		if fmt.Sprintf("%v", i) != string(y) {
			t.Errorf("Encrypt not working")
		}
	}

	for i := 0; i < 100; i++ {
		x, _ := cmhp_crypto.EncryptAes32(fmt.Sprintf("%v", i), fmt.Sprintf("%v", i*10000))
		y, _ := cmhp_crypto.DecryptAes32(x, fmt.Sprintf("%v", i*10000))

		if fmt.Sprintf("%v", i) != string(y) {
			t.Errorf("Encrypt not working")
		}
	}
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmhp_crypto.UID(10)
	}
}
