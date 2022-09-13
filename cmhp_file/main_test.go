package cmhp_file_test

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"os"
	"testing"
)

func TestA(t *testing.T) {
	// l, _ := cmhp_file.ListAll("../")
	// cmhp_print.Print(l)
}

func TestB(t *testing.T) {
	fmt.Printf("%v", cmhp_file.Exists(""))
}

func TestCompress(t *testing.T) {
	x := struct {
		A string
	}{A: "Hix"}
	err := cmhp_file.WriteCompressed(os.TempDir()+"/test", &x)
	if err != nil {
		t.Fatal()
	}

	err = cmhp_file.ReadCompressedJSON(os.TempDir()+"/test", &x)
	if err != nil {
		t.Fatal()
	}

	if x.A != "Hix" {
		t.Fatal()
	}
}
