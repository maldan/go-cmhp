package cmhp_print_test

import (
	"github.com/maldan/go-cmhp/cmhp_print"
	"testing"
	"time"
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

type Person1 struct {
	W3ID string
	Name string
}

type Address1 struct {
	city    string
	country string
}

type User1 struct {
	name      string
	age       int
	address   Address1
	manager   Person1
	developer Person1
	tech      Person1
}

func TestA(t *testing.T) {
	t2 := time.Now()

	cmhp_print.Print(t2)
}

func TestB(t *testing.T) {
	load := User1{
		name: "John Doe",
		age:  34,
		address: Address1{
			city:    "New York",
			country: "USA",
		},
		manager: Person1{
			W3ID: "jBult@in.org.com",
			Name: "Bualt",
		},
		developer: Person1{
			W3ID: "tsumi@in.org.com",
			Name: "Sumi",
		},
		tech: Person1{
			W3ID: "lPaul@in.org.com",
			Name: "Paul",
		},
	}

	cmhp_print.Print(load)
}

func TestC(t *testing.T) {
	cmhp_print.Print(nil)
}
