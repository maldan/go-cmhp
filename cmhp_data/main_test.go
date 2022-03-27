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

func TestIsEnd(t *testing.T) {
	ba := cmhp_data.Allocate(0, true)
	ba.WriteUInt32(32)
	ba.Position = 0
	if ba.IsEnd() {
		t.Errorf("Incorrect value")
	}
	ba.ReadUint32()
	if !ba.IsEnd() {
		t.Errorf("Incorrect value")
	}

	ba.Position = 0
	ba.ReadUint8()
	ba.ReadUint8()
	ba.ReadUint8()
	if ba.IsEnd() {
		t.Errorf("Incorrect value")
	}
	ba.ReadUint8()
	if !ba.IsEnd() {
		t.Errorf("Incorrect value")
	}
}

func TestFloat32(t *testing.T) {
	ba := cmhp_data.Allocate(0, true)
	ba.WriteFloat32(32.32)
	ba.Position = 0
	if ba.ReadFloat32() != 32.32 {
		t.Errorf("Incorrect value")
	}
}

func TestSection(t *testing.T) {
	ba := cmhp_data.Allocate(0, true)

	a := cmhp_data.Allocate(4, true)
	a.WriteUInt32(1)
	b := cmhp_data.Allocate(4, true)
	b.WriteUInt32(1)

	ba.WriteSection(1234, "X", a)
	ba.WriteSection(1234, "Y", b)
	ba.Position = 0

	s1, ss1, _ := ba.ReadSection(1234)
	s2, ss2, _ := ba.ReadSection(1234)

	if s1 != "X" {
		t.Errorf("Incorrect value")
	}
	if s2 != "Y" {
		t.Errorf("Incorrect value")
	}

	if ss1.Length != 4 {
		t.Errorf("Incorrect value")
	}
	if ss2.Length != 4 {
		t.Errorf("Incorrect value")
	}
}
