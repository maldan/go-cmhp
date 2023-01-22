package cmhp_byte_test

import (
	"github.com/maldan/go-cmhp/cmhp_byte"
	"reflect"
	"testing"
)

type TestString struct {
	StrDefault string
	StrTinyRus string `len:"4"`
	StrTinyJp  string `len:"6"`
	StrTiny    string `len:"32"`
	StrShort   string `len:"2048"`
	StrMed     string `len:"1048576"`
}

func loop(size int) string {
	out := make([]byte, size)
	for i := 0; i < size; i++ {
		out[i] = byte(i)
	}

	return string(out)
}

func TestPackCheckBytes(t *testing.T) {
	bytes := cmhp_byte.Pack[TestString](&TestString{StrDefault: "1", StrTiny: "123", StrShort: "123456"})

	// Num of fields
	if bytes[0] != byte(reflect.TypeOf(TestString{}).NumField()) {
		t.Error("Fuck")
	}

	// Length of field
	if bytes[1] != byte(10) {
		t.Error("Fuck")
	}

	// Length of field
	if string(bytes[2:12]) != "StrDefault" {
		t.Error("Fuck")
	}
}

func TestPackString(t *testing.T) {
	bytes := cmhp_byte.Pack[TestString](&TestString{StrDefault: "1", StrTiny: "123", StrShort: "123456"})
	tt := cmhp_byte.Unpack[TestString](&bytes)
	if tt.StrDefault != "1" {
		t.Error("Fuck")
	}
	if tt.StrTiny != "123" {
		t.Error("Fuck")
	}
	if tt.StrShort != "123456" {
		t.Error("Fuck")
	}
}

func TestPackStringMaxLength(t *testing.T) {
	sIn := &TestString{
		StrDefault: loop(65535 * 2),
		StrTiny:    loop(1024),
		StrShort:   loop(4096),
		StrMed:     loop(1_048_576 * 2),
	}

	// Pack and unpack
	bytes := cmhp_byte.Pack[TestString](sIn)
	sOut := cmhp_byte.Unpack[TestString](&bytes)

	if len(sOut.StrDefault) != 65535 {
		t.Errorf("Fuck %v", len(sOut.StrDefault))
	}

	if len(sOut.StrTiny) != 32 {
		t.Error("Fuck")
	}
	if len(sOut.StrShort) != 4096 {
		t.Error("Fuck")
	}
	if len(sOut.StrMed) != 1_048_576 {
		t.Error("Fuck")
	}
}

func TestPackStringMaxLengthUTF8(t *testing.T) {
	sIn := &TestString{StrTinyRus: "сука блядь", StrTinyJp: "系属やから"}

	// Pack and unpack
	bytes := cmhp_byte.Pack[TestString](sIn)
	sOut := cmhp_byte.Unpack[TestString](&bytes)

	if sOut.StrTinyRus != "су" {
		t.Error("Fuck")
	}
	if sOut.StrTinyJp != "系属" {
		t.Error("Fuck")
	}
}
