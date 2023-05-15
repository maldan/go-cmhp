package cmhp_print

import (
	"encoding/binary"
	"fmt"
	"github.com/fatih/color"
	"reflect"
	"strings"
	"unsafe"
)

func printMap(s interface{}, ident int) {
	ident1 := strings.Repeat(" ", ident)
	ident2 := strings.Repeat(" ", ident+2)

	fmt.Print("{\n")
	c := color.New(color.FgGreen)

	sv := reflect.ValueOf(s)
	for n, e := range sv.MapKeys() {
		v := sv.MapIndex(e)
		kn := sv.MapKeys()[n]

		fmt.Printf("%v", ident2)
		c.Printf("%v", kn)
		fmt.Print(" = ")
		__print(v.Interface(), ident+2, 0, true)
	}

	fmt.Print(ident1 + "}")
}

func printStruct(s interface{}, ident int, arrayIdent int, isHideType bool) {
	identArr := strings.Repeat(" ", arrayIdent)
	ident1 := strings.Repeat(" ", ident)
	ident2 := strings.Repeat(" ", ident+2)

	sv := reflect.ValueOf(s)
	st := reflect.TypeOf(s)
	sv2 := reflect.New(sv.Type()).Elem()
	sv2.Set(sv)

	// Print struct header
	fmt.Print(identArr)

	c := color.New(color.FgRed).Add(color.Bold)
	c1 := color.New(color.FgCyan).Add(color.Underline)

	if !isHideType {
		c.Print("struct ")
		c1.Print(sv.Type())
		fmt.Print(" ")
	}

	fmt.Print("{\n")

	// Print struct fields
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)

		rf := sv2.Field(i)
		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()

		fmt.Printf("%v", ident2)
		c.Print(sv.Field(i).Type())
		fmt.Printf(" %v = ", f.Name)

		__print(rf.Interface(), ident+2, 0, true)
	}

	fmt.Print(ident1 + "}")
}

func printSlice(s interface{}, ident int) {
	ident1 := strings.Repeat(" ", ident)

	sv := reflect.ValueOf(s)

	// Print slice header
	//c1 := color.New(color.FgCyan).Add(color.Underline)
	//c1.Print(sv.Type())
	/*if !isHeaderOnSameLine {
		fmt.Print(ident1)
	}*/

	fmt.Print("[\n")

	for i := 0; i < sv.Len(); i++ {
		__print(sv.Index(i).Interface(), ident+2, ident+2, false)
	}
	fmt.Printf("%v]", ident1)
}

func __print(val any, ident int, arrayIdent int, isHideType bool) {
	v := reflect.TypeOf(val)
	if v == nil {
		fmt.Printf("%v\n", val)
		return
	}
	if v.Kind() == reflect.Pointer {
		el := reflect.ValueOf(val)
		if el.IsNil() {
			__print(nil, ident, arrayIdent, isHideType)
		} else {
			__print(reflect.ValueOf(val).Elem().Interface(), ident, arrayIdent, isHideType)
		}
	} else if v.Kind() == reflect.Map {
		printMap(val, ident)
	} else if v.Kind() == reflect.Struct {
		printStruct(val, ident, arrayIdent, isHideType)
	} else if v.Kind() == reflect.Slice {
		printSlice(val, ident)
	} else if v.Kind() == reflect.String {
		c := color.New(color.FgGreen)
		c.Print(val)
	} else if v.Kind() == reflect.Int ||
		v.Kind() == reflect.Int8 || v.Kind() == reflect.Int16 || v.Kind() == reflect.Int32 || v.Kind() == reflect.Int64 ||
		v.Kind() == reflect.Uint8 || v.Kind() == reflect.Uint16 || v.Kind() == reflect.Uint32 || v.Kind() == reflect.Uint64 ||
		v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64 {
		c := color.New(color.FgBlue)
		c.Print(val)
	} else {
		fmt.Printf("%v", val)
	}
	fmt.Print("\n")
}

func Print(val interface{}) {
	__print(val, 0, 0, false)
}

const BgRed = color.BgRed
const BgGreen = color.BgGreen
const BgBlue = color.BgBlue

type ColorRange struct {
	From  int
	Len   int
	Color color.Attribute
}

func PrintBytesColored(b []byte, lineSize int, colors []ColorRange) {
	buff := make([]string, 0)
	for i := 0; i < len(b); i++ {
		c := color.New()
		// c := color.New(color.FgRed)
		for j := 0; j < len(colors); j++ {
			if i >= colors[j].From && i < colors[j].From+colors[j].Len {
				c = color.New(colors[j].Color)
			}
		}

		buff = append(buff, c.Sprintf("%02X ", b[i]))
		if i != 0 && (i+1)%lineSize == 0 {
			buff = append(buff, "\n")
		}
	}
	buff = append(buff, "\n")
	fmt.Print(strings.Join(buff, ""))
}

func PrintBytes(b []byte, lineSize int) {
	buff := make([]string, 0)
	for i := 0; i < len(b); i++ {
		buff = append(buff, fmt.Sprintf("%02X ", b[i]))
		if i != 0 && (i+1)%lineSize == 0 {
			buff = append(buff, "\n")
		}
	}
	buff = append(buff, "\n")
	fmt.Print(strings.Join(buff, ""))
}

func PrintDebugBytes(b []byte, lineBreaks ...int) {
	o := 0
	c := 0
	for i := 0; i < len(b); i++ {
		c += 1
		fmt.Printf("%02X ", b[i])
		if o < len(lineBreaks) && lineBreaks[o] == c {
			if lineBreaks[o] == 4 {
				fmt.Printf(" - %v", binary.LittleEndian.Uint32(b[i-3:]))
			}
			fmt.Printf("\n")
			o += 1
			c = 0
		}

		/*if lineBreaks[i] {
			fmt.Printf("\n")
		}*/
	}
	// fmt.Printf("%x", b)
}
