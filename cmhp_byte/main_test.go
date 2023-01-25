package cmhp_byte_test

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_byte"
	"testing"
)

type Test struct {
	Id  uint8  `json:"id"`
	Id2 uint16 `json:"id2"`
	Id3 uint32 `json:"id3"`

	S string `json:"s"`

	FirstName   string `json:"first_name" len:"32"`
	LastName    string `json:"last_name" len:"2048"`
	DateOfBirth string `json:"date_of_birth"`

	Age        uint32 `json:"age"`
	UserId     uint32 `json:"user_id"`
	DocumentId uint32 `json:"document_id"`
	ReviewId   uint32 `json:"review_id"`

	Address               string `json:"address"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	ZipCode               string `json:"zip_code"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	EmergencyPersonName   string `json:"emergency_person_name"`
	EmergencyPersonPhone  string `json:"emergency_person_phone"`
	ParentOrLegalGuardian string `json:"parent_or_legal_guardian"`
	Date                  string `json:"date"`

	// Arr [4]int `json:"arr"`

	/*
		ReviewEnabled         bool   `json:"review_enabled"`

		Created time.Time `json:"created"`*/
}

func TestUintX(t *testing.T) {
	u := uint8(255)
	i := int8(u)
	fmt.Printf("%v\n", u)
	fmt.Printf("%v\n", i)
	buff := make([]byte, 4)

	cmhp_byte.From8ToBuffer(&u, &buff, 0)
	cmhp_byte.From8ToBuffer(&i, &buff, 1)

	u2 := uint8(0)
	i2 := int8(0)

	cmhp_byte.FromBufferTo8(&buff, &u2, 0)
	cmhp_byte.FromBufferTo8(&buff, &i2, 1)

	fmt.Printf("%v\n", u2)
	fmt.Printf("%v\n", i2)
}

func TestUint8(t *testing.T) {
	u := uint8(234)
	buff := make([]byte, 4)

	cmhp_byte.From8ToBuffer(&u, &buff, 0)
	if buff[0] != 234 {
		t.Errorf("Not working")
	}

	cmhp_byte.From8ToBuffer(&u, &buff, 1)
	if buff[1] != 234 {
		t.Errorf("Not working")
	}

	u2 := byte(96)
	buff2 := make([]byte, 4)

	cmhp_byte.From8ToBuffer(&u2, &buff2, 0)
	if buff2[0] != 96 {
		t.Errorf("Not working")
	}

	cmhp_byte.From8ToBuffer(&u2, &buff2, 1)
	if buff2[1] != 96 {
		t.Errorf("Not working")
	}
}

func TestUint16(t *testing.T) {
	u := uint16(1024)
	buff := make([]byte, 4)

	cmhp_byte.From16ToBuffer(&u, &buff, 0)

	if buff[0] != 0 || buff[1] != 4 {
		t.Errorf("Not working")
	}
}

func TestUint64(t *testing.T) {
	u := uint64(1)
	buff := make([]byte, 8)

	cmhp_byte.Write64ToBuffer(u, &buff, 0)

	fmt.Printf("%v\n", buff)
}

func TestBuffer(t *testing.T) {
	u := uint8(146)
	buff := make([]byte, 4)
	fmt.Printf("%v\n", u)

	cmhp_byte.From8ToBuffer(&u, &buff, 0)
	fmt.Printf("%v\n", u)

	u2 := uint8(0)
	cmhp_byte.FromBufferTo8(&buff, &u2, 0)

	if u2 != u {
		t.Errorf("Not working")
	}
}

func TestPack(t *testing.T) {
	bytes := cmhp_byte.Pack(&Test{Id: 64, Id2: 128, Id3: 1048576, S: "Hi Lol",
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	})
	fmt.Printf("%v\n", bytes)
	fmt.Printf("%v\n", string(bytes))
}

func TestUnpack(t *testing.T) {
	bytes := cmhp_byte.Pack(&Test{Id: 64, Id2: 128, Id3: 1048576, S: "Hi Lol",
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	})
	tt := cmhp_byte.Unpack[Test](bytes)
	fmt.Printf("%v\n", tt)
}

func BenchmarkMy2(b *testing.B) {
	a := Test{
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	}
	for i := 0; i < b.N; i++ {
		v := cmhp_byte.Pack(&a)
		b.SetBytes(int64(len(v)))
	}
}

func BenchmarkUnMy(b *testing.B) {
	a := Test{
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	}
	bytes := cmhp_byte.Pack(&a)

	for i := 0; i < b.N; i++ {
		cmhp_byte.Unpack[Test](bytes)
		b.SetBytes(int64(len(bytes)))
	}
}
