package cmhp_file_test

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"testing"
)

func TestA(t *testing.T) {
	// l, _ := cmhp_file.ListAll("../")
	// cmhp_print.Print(l)
}

func TestB(t *testing.T) {
	fmt.Printf("%v", cmhp_file.Exists(""))
}
