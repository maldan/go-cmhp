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

func TestAllow(t *testing.T) {
	if cmhp_string.Allow("Hello World", "el") != "elll" {
		t.Fatalf("Not working")
	}
	if cmhp_string.Allow("123", "1") != "1" {
		t.Fatalf("Not working")
	}
	if cmhp_string.Allow("Keks, city", "Keks city") != "Keks city" {
		t.Fatalf("Not working")
	}
	if cmhp_string.AllowCommon("Keks, city") != "Keks city" {
		t.Fatalf("Not working")
	}
}

func TestEmail(t *testing.T) {
	if cmhp_string.IsEmailValid("hello@world.ru") != true {
		t.Fatalf("Not working")
	}
	if cmhp_string.IsEmailValid("hello@ world.ru") != false {
		t.Fatalf("Not working")
	}
	if cmhp_string.IsEmailValid(" hello@world.ru") != false {
		t.Fatalf("Not working")
	}
	if cmhp_string.IsEmailValid("hello@world.ru ") != false {
		t.Fatalf("Not working")
	}
}

func BenchmarkA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmhp_string.LowerFirst("SasageoDavageo")
	}
}
