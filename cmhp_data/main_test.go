package cmhp_data_test

import (
	"github.com/maldan/go-cmhp/cmhp_data"
	"testing"
)

func TestAllocate(t *testing.T) {
	ba := cmhp_data.Allocate(0, true)
	ba.WriteUInt16(1024)
	ba.WriteUInt32(1024)
	ba.WriteUInt32(1024)

	if ba.Length != 10 {
		t.Errorf("Incorrect allocation")
	}
}

func TestLength(t *testing.T) {
	ba := cmhp_data.Allocate(4, true)
	ba.WriteUInt32(1024)
	ba.Position = 0
	ba.WriteUInt32(1024)

	if ba.Length != 4 {
		t.Errorf("Incorrect length")
	}
}

func TestWriteAndReadUint8(t *testing.T) {
	ba := cmhp_data.Allocate(4, true)
	ba.WriteUInt16(1024)
	ba.Position = 0
	x := ba.ReadUint16()

	if x != 1024 {
		t.Errorf("Incorrect value")
	}
}

func TestWriteAndReadString(t *testing.T) {
	ba := cmhp_data.Allocate(0, true)
	ba.WriteUTF8("Hello")
	ba.WriteUTF8("Fuck you")
	ba.WriteUTF8("Сасагео")
	ba.Position = 0

	if ba.ReadUTF8() != "Hello" {
		t.Errorf("Incorrect value")
	}
	if ba.ReadUTF8() != "Fuck you" {
		t.Errorf("Incorrect value")
	}
	if ba.ReadUTF8() != "Сасагео" {
		t.Errorf("Incorrect value")
	}
}
