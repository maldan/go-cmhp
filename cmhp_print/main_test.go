package cmhp_print_test

import (
	"github.com/maldan/go-cmhp/cmhp_print"
	"testing"
)

type D struct {
	O int
	R int
}

type C struct {
	A int
	B int
	G []D
}

type B struct {
	Z int
	C C
}

type A struct {
	X  int
	Y  int
	IN B
}

func TestA(t *testing.T) {

	a := A{}

	cmhp_print.Print(a)
}
