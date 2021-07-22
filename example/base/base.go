package main

import (
	"fmt"

	"github.com/maldan/go-cmhp"
)

type Lox struct {
	X int
}

func main() {
	cmhp.FileWriteAsText("a/b/sas.txt", "sasgeoo")
	cmhp.FileWriteAsText("a/b/sa2s.txt", "sasgeoo")

	lox := make([]interface{}, 0)
	for i := 0; i < 10; i++ {
		lox = append(lox, Lox{X: i})
	}

	x := cmhp.SliceFilter(lox, func(i interface{}) bool {
		return i.(Lox).X > 2
	})
	fmt.Println(x)

	y := cmhp.SliceMap(lox, func(i interface{}) interface{} {
		return i.(Lox).X > 5
	})
	fmt.Println(y)

	z := cmhp.SliceSort(lox, func(i, j int) bool {
		return lox[i].(Lox).X > lox[j].(Lox).X
	})
	fmt.Println(z)

	// fmt.Println(cmhp.RequestGetAsText("https://api.github.com/repos/maldan/gamx/releases", nil))
}
