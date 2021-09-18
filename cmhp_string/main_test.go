package cmhp_string_test

import (
	"testing"

	"github.com/maldan/go-cmhp/cmhp_string"
)

func TestA(t *testing.T) {
	if cmhp_string.LowerFirst("SasageoDavageo") != "sasageoDavageo" {
		t.Fatalf("Not working")
	}
}

func BenchmarkA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmhp_string.LowerFirst("SasageoDavageo")
	}
}
