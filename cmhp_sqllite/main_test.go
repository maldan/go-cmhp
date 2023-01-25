package cmhp_sqllite_test

import (
	"github.com/maldan/go-cmhp/cmhp_sqllite"
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
	/*r := cmhp_sqllite.CreateTable[Test]("sas")
	fmt.Printf("%v\n", r)*/

	cmhp_sqllite.Insert(nil, "sas", Test{
		Id:                   10,
		Id3:                  20,
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
}
