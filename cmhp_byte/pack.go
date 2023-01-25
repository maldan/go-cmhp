package cmhp_byte

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

type (
	TpInfo struct {
		HeaderSize     int
		Name           [][]byte
		Offset         []uintptr
		Type           []uint8
		MaxBytesLength []int
		FieldAmount    uint8

		MapOffset map[string]uintptr
		MapType   map[string]uint8
	}
)

var typeCache = map[reflect.Type]*TpInfo{}

func getTypeInfo[T any](v *T) *TpInfo {
	vv, ok := typeCache[reflect.TypeOf(v).Elem()]
	if ok {
		return vv
	}

	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()
	info := TpInfo{}

	info.HeaderSize = 1 // num of fields
	info.MapOffset = map[string]uintptr{}
	info.MapType = map[string]uint8{}
	info.MaxBytesLength = make([]int, typeOf.NumField())
	info.Offset = make([]uintptr, typeOf.NumField())
	info.Name = make([][]byte, typeOf.NumField())
	info.Type = make([]uint8, typeOf.NumField())

	for i := 0; i < typeOf.NumField(); i++ {
		info.Offset[i] = typeOf.Field(i).Offset
		info.MapOffset[typeOf.Field(i).Name] = typeOf.Field(i).Offset
		info.Name[i] = []byte(typeOf.Field(i).Name)

		// Set header size
		info.HeaderSize += 1                         // field len
		info.HeaderSize += len(typeOf.Field(i).Name) // len
		info.HeaderSize += 1                         // type
		info.FieldAmount += 1

		switch valueOf.Field(i).Interface().(type) {
		case bool:
			info.Type[i] = _bool
			info.MapType[typeOf.Field(i).Name] = _bool
			info.HeaderSize += 1
			break
		case uint8, int8:
			info.Type[i] = _i8
			info.MapType[typeOf.Field(i).Name] = _i8
			info.HeaderSize += 1
			break
		case uint16, int16:
			info.Type[i] = _i16
			info.MapType[typeOf.Field(i).Name] = _i16
			info.HeaderSize += 2
			break
		case uint32, int32, int:
			info.Type[i] = _i32
			info.MapType[typeOf.Field(i).Name] = _i32
			info.HeaderSize += 4
			break
		case float32:
			info.Type[i] = _f32
			info.MapType[typeOf.Field(i).Name] = _f32
			info.HeaderSize += 4
			break
		case uint64, int64:
			info.Type[i] = _i64
			info.MapType[typeOf.Field(i).Name] = _i64
			info.HeaderSize += 8
			break
		case float64:
			info.Type[i] = _f64
			info.MapType[typeOf.Field(i).Name] = _f64
			info.HeaderSize += 8
			break
		case string:
			n, _ := strconv.Atoi(typeOf.Field(i).Tag.Get("len"))
			info.MaxBytesLength[i] = n

			if n == 0 { // default max length up to 64 kb
				info.Type[i] = _stringShort
				info.MapType[typeOf.Field(i).Name] = _stringShort
				info.HeaderSize += 2 // length of string
				info.MaxBytesLength[i] = 65535
				break
			}
			if n <= 255 {
				info.Type[i] = _stringTiny
				info.MapType[typeOf.Field(i).Name] = _stringTiny
				info.HeaderSize += 1 // length of string
				break
			} else if n <= 65535 {
				info.Type[i] = _stringShort
				info.MapType[typeOf.Field(i).Name] = _stringShort
				info.HeaderSize += 2 // length of string
				break
			} else if n <= 16_777_215 {
				info.Type[i] = _stringMed
				info.MapType[typeOf.Field(i).Name] = _stringMed
				info.HeaderSize += 3 // length of string
				break
			} else {
				info.Type[i] = _stringBig
				info.MapType[typeOf.Field(i).Name] = _stringBig
				info.HeaderSize += 4 // length of string
				break
			}

		default:
			panic(fmt.Sprintf("unsupported type %T", valueOf.Field(i).Interface()))
		}
	}

	typeCache[typeOf] = &info
	return &info
}

func Pack[T any](v *T) []byte {
	info := getTypeInfo(v)

	// pos in array
	p := 0

	// Calculate string size
	size := info.HeaderSize
	for i := 0; i < len(info.Type); i++ {
		if info.Type[i] == _stringTiny || info.Type[i] == _stringShort || info.Type[i] == _stringMed || info.Type[i] == _stringBig {
			bb := *(*[]byte)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))
			ln := len(bb)
			size += ln
		}
	}

	// Prepare buffer
	out := make([]byte, size)

	// Header
	out[p] = info.FieldAmount
	p += 1

	// Content
	for i := 0; i < len(info.Offset); i++ {
		// Field name
		p += FromString8ToBuffer(&info.Name[i], &out, p)

		// Type info
		out[p] = info.Type[i]
		p += 1

		switch info.Type[i] {
		case _bool:
			bl := *(*bool)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))
			if bl {
				out[p] = 1
			} else {
				out[p] = 0
			}
			p += 1
			break
		case _i8:
			out[p] = *(*uint8)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))
			p += 1
			break
		case _i16:
			rv := *(*uint16)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))
			p += Write16ToBuffer(rv, &out, p)
			break
		case _i32:
			rv := *(*uint32)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))
			p += Write32ToBuffer(rv, &out, p)
			break
		case _i64:
			rv := *(*uint64)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))
			p += Write64ToBuffer(rv, &out, p)
			break
		case _stringTiny:
			rv := (*string)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))

			// Check max string length
			ln := len(*rv)
			if ln > info.MaxBytesLength[i] {
				ln = info.MaxBytesLength[i]
			}
			rv2 := (*rv)[:ln]

			p += FromString8ToBuffer(&rv2, &out, p)
			break
		case _stringShort:
			rv := (*string)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))

			// Check max string length
			/*ln := len(*rv)
			if ln > info.MaxBytesLength[i] {
				ln = info.MaxBytesLength[i]
			}
			// fmt.Printf("WRITE AS: %v - %v\n", string(info.Name[i]), ln)
			rv2 := (*rv)[:ln]
			// fmt.Printf("%v - %v\n", p, []byte(rv2))*/

			p += FromString16ToBuffer(rv, &out, p)
			break
		case _stringMed:
			rv := (*string)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))

			// Check max string length
			ln := len(*rv)
			if ln > info.MaxBytesLength[i] {
				ln = info.MaxBytesLength[i]
			}
			rv2 := (*rv)[:ln]

			p += FromString24ToBuffer(&rv2, &out, p)
			break
		case _stringBig:
			rv := (*string)(unsafe.Add(unsafe.Pointer(v), info.Offset[i]))

			// Check max string length
			ln := len(*rv)
			if ln > info.MaxBytesLength[i] {
				ln = info.MaxBytesLength[i]
			}
			rv2 := (*rv)[:ln]

			p += FromString32ToBuffer(&rv2, &out, p)
			break
		}
	}

	// 48 nc - base

	return out
}

func Unpack[T any](b []byte) T {
	out := new(T)
	info := getTypeInfo(out)

	// Read fields number
	p := 0
	numOfFields := (b)[p]
	p++

	// Read fields
	for i := 0; i < int(numOfFields); i++ {
		// Read field length
		fieldLength := int((b)[p])
		p += 1

		// Read field name
		fieldName := string((b)[p : p+fieldLength])
		p += fieldLength

		// Read field type
		vType := (b)[p]
		p += 1

		// Field offset
		offset, _ := info.MapOffset[fieldName]

		switch vType {
		case _i8:
			*(*uint8)(unsafe.Add(unsafe.Pointer(out), offset)) = (b)[p]
			p += 1
			break
		case _i16:
			p += FromBufferTo16(&b, (*uint16)(unsafe.Add(unsafe.Pointer(out), offset)), p)
			break
		case _i32:
			p += FromBufferTo32(&b, (*uint32)(unsafe.Add(unsafe.Pointer(out), offset)), p)
			break
		case _stringTiny:
			// string length
			strLength := (b)[p]
			p += 1
			//*(*string)(unsafe.Add(unsafe.Pointer(out), offset)) = string((b)[p : p+int(strLength)])
			*(*[]byte)(unsafe.Add(unsafe.Pointer(out), offset)) = make([]byte, int(strLength))
			copy(*(*[]byte)(unsafe.Add(unsafe.Pointer(out), offset)), (b)[p:p+int(strLength)])
			// fmt.Printf("%v", "XX")
			p += int(strLength)
			break
		case _stringShort:
			// fmt.Printf("%v\n", "GAS")
			// string length
			strLength := Read16FromBuffer(b, p)
			p += 2
			*(*string)(unsafe.Add(unsafe.Pointer(out), offset)) = string((b)[p : p+int(strLength)])
			p += int(strLength)
			break
		case _stringMed:
			// string length
			strLength := Read24FromBuffer(b, p)
			p += 3
			*(*string)(unsafe.Add(unsafe.Pointer(out), offset)) = string((b)[p : p+int(strLength)])
			p += int(strLength)
			break
		case _stringBig:
			// string length
			strLength := Read32FromBuffer(b, p)
			p += 4
			*(*string)(unsafe.Add(unsafe.Pointer(out), offset)) = string((b)[p : p+int(strLength)])
			p += int(strLength)
			break
		}
	}

	return *out
}
