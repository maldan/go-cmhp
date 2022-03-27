package cmhp_data

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"math"
)

type ByteArray struct {
	Position       uint64
	Capacity       uint64
	Length         uint64
	Data           []byte
	IsLittleEndian bool
}

func Allocate(size uint64, isLE bool) *ByteArray {
	ba := ByteArray{Capacity: size, Data: make([]byte, size), IsLittleEndian: isLE}
	return &ba
}

func FromFile(path string, isLE bool) (*ByteArray, error) {
	data, err := cmhp_file.ReadBin(path)
	if err != nil {
		return nil, err
	}
	ba := ByteArray{Capacity: uint64(len(data)), Length: uint64(len(data)), Data: data, IsLittleEndian: isLE}
	return &ba, nil
}

func (b *ByteArray) WriteUInt8(value uint8) {
	if b.Length+1 > b.Capacity {
		b.Grow(1)
	}
	b.Data[b.Position] = value
	b.Position += 1
	if b.Position > b.Length {
		b.Length += 1
	}
}

func (b *ByteArray) WriteUInt16(value uint16) {
	if b.IsLittleEndian {
		b.WriteUInt8(uint8(value & 0xff))
		b.WriteUInt8(uint8(value >> 8))
	} else {
		b.WriteUInt8(uint8(value >> 8))
		b.WriteUInt8(uint8(value & 0xff))
	}
}

func (b *ByteArray) WriteUInt32(value uint32) {
	if b.IsLittleEndian {
		b.WriteUInt8(uint8(value & 0xff))
		b.WriteUInt8(uint8((value >> 8) & 0xff))
		b.WriteUInt8(uint8((value >> 16) & 0xff))
		b.WriteUInt8(uint8(value >> 24))
	} else {
		b.WriteUInt8(uint8(value >> 24))
		b.WriteUInt8(uint8((value >> 16) & 0xff))
		b.WriteUInt8(uint8((value >> 8) & 0xff))
		b.WriteUInt8(uint8(value & 0xff))
	}
}

func (b *ByteArray) WriteInt8(value int8) {
	b.WriteUInt8(uint8(value))
}

func (b *ByteArray) WriteInt16(value int16) {
	b.WriteUInt16(uint16(value))
}

func (b *ByteArray) WriteInt32(value int32) {
	b.WriteUInt32(uint32(value))
}

func (b *ByteArray) WriteUTF8(value string) {
	b.WriteUInt32(uint32(len(value)))
	arr := []byte(value)
	for i := 0; i < len(arr); i++ {
		b.WriteUInt8(arr[i])
	}
}

func (b *ByteArray) WriteFloat32(value float32) {
	n := math.Float32bits(value)

	if b.IsLittleEndian {
		b.WriteUInt8(uint8(n & 0xff))
		b.WriteUInt8(uint8((n >> 8) & 0xff))
		b.WriteUInt8(uint8((n >> 16) & 0xff))
		b.WriteUInt8(uint8(n >> 24))
	} else {
		b.WriteUInt8(uint8(n >> 24))
		b.WriteUInt8(uint8((n >> 16) & 0xff))
		b.WriteUInt8(uint8((n >> 8) & 0xff))
		b.WriteUInt8(uint8(n & 0xff))
	}
}

func (b *ByteArray) WriteSection(marker uint32, name string, arr *ByteArray) {
	b.WriteUInt32(marker)
	b.WriteUTF8(name)

	// Data with length
	b.WriteUInt32(uint32(arr.Length))
	for i := 0; i < int(arr.Length); i++ {
		b.WriteUInt8(arr.Data[i])
	}
}

func (b *ByteArray) ReadFloat32() float32 {
	bytes := make([]byte, 4)
	bytes[0] = b.ReadUint8()
	bytes[1] = b.ReadUint8()
	bytes[2] = b.ReadUint8()
	bytes[3] = b.ReadUint8()
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func (b *ByteArray) ReadUint8() uint8 {
	b.Position += 1
	return b.Data[b.Position-1]
}

func (b *ByteArray) ReadUint16() uint16 {
	b1 := b.ReadUint8()
	b2 := b.ReadUint8()

	if b.IsLittleEndian {
		return uint16(int(b1) + int(b2)*256)
	} else {
		return uint16(int(b1)*256 + int(b2))
	}
}

func (b *ByteArray) ReadUint32() uint32 {
	b1 := b.ReadUint8()
	b2 := b.ReadUint8()
	b3 := b.ReadUint8()
	b4 := b.ReadUint8()

	if b.IsLittleEndian {
		return uint32(int(b1) + int(b2)*256 + int(b3)*65536 + int(b4)*16777216)
	} else {
		return uint32(int(b1)*16777216 + int(b2)*65536 + int(b3)*256 + int(b4))
	}
}

func (b *ByteArray) ReadUTF8() string {
	l := b.ReadUint32()
	str := b.Read(uint64(l))
	return string(str)
}

func (b *ByteArray) Read(amount uint64) []byte {
	data := b.Data[b.Position : b.Position+amount]
	b.Position += amount
	return data
}

func (b *ByteArray) ReadSection(marker uint32) (string, *ByteArray, error) {
	// Read header
	m := b.ReadUint32()
	if m != marker {
		return "", nil, errors.New("unknown marker symbol")
	}
	name := b.ReadUTF8()
	size := b.ReadUint32()

	// Read data
	b2 := Allocate(0, b.IsLittleEndian)
	for i := 0; i < int(size); i++ {
		b2.WriteUInt8(b.ReadUint8())
	}
	b2.Position = 0

	return name, b2, nil
}

func (b *ByteArray) Grow(amount uint64) {
	newArray := make([]byte, b.Capacity+amount)
	for i := 0; i < int(b.Length); i++ {
		newArray[i] = b.Data[i]
	}
	b.Data = newArray
	b.Capacity += amount
}

func (b *ByteArray) Print() {
	for i := 0; i < int(b.Length); i++ {
		fmt.Printf("%02X ", b.Data[i])
	}
	fmt.Printf("\n")
}

func (b *ByteArray) IsEnd() bool {
	return b.Position >= b.Length
}
