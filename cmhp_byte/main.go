package cmhp_byte

import (
	"encoding/binary"
	"unsafe"
)

func From8ToBuffer[T uint8 | int8](from *T, buffer *[]byte, bufferOffset int) int {
	hh := *(*[1]byte)(unsafe.Pointer(from))
	copy((*buffer)[bufferOffset:], hh[:])
	return 1
}

func From16ToBuffer[T uint16 | int16](from *T, buffer *[]byte, bufferOffset int) int {
	hh := *(*[2]byte)(unsafe.Pointer(from))
	copy((*buffer)[bufferOffset:], hh[:])
	return 2
}

func From24ToBuffer[T uint32 | int32](from *T, buffer *[]byte, bufferOffset int) int {
	hh := *(*[3]byte)(unsafe.Pointer(from))
	copy((*buffer)[bufferOffset:], hh[:])
	return 3
}

func From32ToBuffer[T uint32 | int32](from *T, buffer *[]byte, bufferOffset int) int {
	hh := *(*[4]byte)(unsafe.Pointer(from))
	copy((*buffer)[bufferOffset:], hh[:])
	return 4
}

func From64ToBuffer[T uint64 | int64](from *T, buffer *[]byte, bufferOffset int) int {
	hh := *(*[8]byte)(unsafe.Pointer(from))
	copy((*buffer)[bufferOffset:], hh[:])
	return 8
}

func FromString8ToBuffer[T string | []byte](from *T, buffer *[]byte, bufferOffset int) int {
	bb := *(*[]byte)(unsafe.Pointer(from))
	bb2 := bb[:]
	rln := len(bb)
	if rln > 255 {
		rln = 255
		bb2 = bb[:255]
	}

	Write8ToBuffer(uint8(rln), buffer, bufferOffset)
	copy((*buffer)[bufferOffset+1:], bb2)
	return 1 + rln
}

func FromString16ToBuffer[T string | []byte](from *T, buffer *[]byte, bufferOffset int) int {
	bb := *(*[]byte)(unsafe.Pointer(from))
	bb2 := bb[:]
	rln := len(bb)
	if rln > 65535 {
		rln = 65535
		bb2 = bb[:65535]
	}

	Write16ToBuffer(uint16(rln), buffer, bufferOffset)
	copy((*buffer)[bufferOffset+2:], bb2)

	return 2 + rln
}

func FromString24ToBuffer[T string | []byte](from *T, buffer *[]byte, bufferOffset int) int {
	bb := *(*[]byte)(unsafe.Pointer(from))
	bb2 := bb[:]
	rln := len(bb)
	if rln > 16_777_215 {
		rln = 16_777_215
		bb2 = bb[:16_777_215]
	}

	Write24ToBuffer(uint32(rln), buffer, bufferOffset)
	copy((*buffer)[bufferOffset+3:], bb2)
	return 3 + rln
}

func FromString32ToBuffer[T string | []byte](from *T, buffer *[]byte, bufferOffset int) int {
	bb := *(*[]byte)(unsafe.Pointer(from))
	ln := len(bb)
	Write32ToBuffer(uint32(ln), buffer, bufferOffset)
	copy((*buffer)[bufferOffset+4:], bb)
	return 4 + ln
}

func FromBufferTo8[T uint8 | int8](buffer *[]byte, dest *T, bufferOffset int) {
	_t := (*[1]byte)(unsafe.Pointer(dest))
	copy(_t[:], (*buffer)[bufferOffset:])
}

func FromBufferTo16[T uint16 | int16](buffer *[]byte, dest *T, bufferOffset int) int {
	_t := (*[2]byte)(unsafe.Pointer(dest))
	copy(_t[:], (*buffer)[bufferOffset:])
	return 2
}

func FromBufferTo32[T uint32 | int32](buffer *[]byte, dest *T, bufferOffset int) int {
	_t := (*[4]byte)(unsafe.Pointer(dest))
	copy(_t[:], (*buffer)[bufferOffset:])
	return 4
}

func FromBufferTo64[T uint64 | int64](buffer *[]byte, dest *T, bufferOffset int) {
	_t := (*[8]byte)(unsafe.Pointer(dest))
	copy(_t[:], (*buffer)[bufferOffset:])
}

func Write8ToBuffer(v uint8, buffer *[]byte, bufferOffset int) int {
	(*buffer)[bufferOffset] = v
	return 1
}

func Write16ToBuffer[T uint16 | int16](v T, buffer *[]byte, bufferOffset int) int {
	off := (*buffer)[bufferOffset : bufferOffset+2]
	binary.LittleEndian.PutUint16(off, uint16(v))
	return 2
}

func Write24ToBuffer(v uint32, buffer *[]byte, bufferOffset int) int {
	off := (*buffer)[bufferOffset:]
	off[0] = uint8(v & 0xff)
	off[1] = uint8((v >> 8) & 0xff)
	off[2] = uint8((v >> 16) & 0xff)
	return 3
}

func Write32ToBuffer(v uint32, buffer *[]byte, bufferOffset int) int {
	off := (*buffer)[bufferOffset : bufferOffset+4]
	binary.LittleEndian.PutUint32(off, v)
	return 4
}

func Write64ToBuffer(v uint64, buffer *[]byte, bufferOffset int) int {
	off := (*buffer)[bufferOffset:]
	binary.LittleEndian.PutUint64(off, v)
	return 8
}

func Read16FromBuffer(buffer []byte, bufferOffset int) uint16 {
	//return binary.LittleEndian.Uint16((buffer)[bufferOffset:bufferOffset+2])
	return uint16(int((buffer)[bufferOffset]) + int((buffer)[bufferOffset+1])*256)
}

func Read24FromBuffer(buffer []byte, bufferOffset int) uint32 {
	return uint32(int((buffer)[bufferOffset]) + int((buffer)[bufferOffset+1])*256 + int((buffer)[bufferOffset+2])*65536)
}

func Read32FromBuffer(buffer []byte, bufferOffset int) uint32 {
	return uint32(int((buffer)[bufferOffset]) + int((buffer)[bufferOffset+1])*256 + int((buffer)[bufferOffset+2])*65536 + int((buffer)[bufferOffset+3])*16777216)
}

func CheckBitMask[T int | uint | uint64 | int64](v T, mask T) bool {
	return v&mask == mask
}
