package main

import (
	"fmt"

	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_process"
	"github.com/maldan/go-cmhp/cmhp_slice"
)

type Lox struct {
	X int
}

func main() {
	cmhp_file.WriteText("a/b/sas.txt", "sasgeoo")
	cmhp_file.WriteText("a/b/sa2s.txt", "sasgeoo")

	lox := make([]interface{}, 0)
	for i := 0; i < 10; i++ {
		lox = append(lox, Lox{X: i})
	}

	x := cmhp_slice.Filter(lox, func(i interface{}) bool {
		return i.(Lox).X > 2
	})
	fmt.Println(x)

	y := cmhp_slice.Map(lox, func(i interface{}) interface{} {
		return i.(Lox).X > 5
	})
	fmt.Println(y)

	z := cmhp_slice.Sort(lox, func(i, j int) bool {
		return lox[i].(Lox).X > lox[j].(Lox).X
	})
	fmt.Println(z)

	// fmt.Println(cmhp.RequestGetAsText("https://api.github.com/repos/maldan/gamx/releases", nil))

	out, _ := cmhp_process.Exec("gam", "process", "list", "--format=json")
	fmt.Println(out)
}
