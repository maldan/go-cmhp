package cmhp_convert_test

import (
	"testing"

	"github.com/maldan/go-cmhp/cmhp_convert"
)

func TestOne(t *testing.T) {
	if cmhp_convert.StrToInt("a") != 0 {
		t.Error("Error")
	}
	if cmhp_convert.StrToInt("5") != 5 {
		t.Error("Error")
	}
}

func BenchmarkOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmhp_convert.StrToInt("5")
	}
}

func BenchmarkTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmhp_convert.IntToStr(i)
	}
}
