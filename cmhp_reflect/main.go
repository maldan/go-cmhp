package cmhp_reflect

import (
	"encoding/json"
	"reflect"
)

func SetField(s interface{}, name string, v interface{}) {
	f := reflect.ValueOf(s).Elem().FieldByName(name)
	if f.CanSet() {
		f.Set(reflect.ValueOf(v))
	}
}

func GetFieldValue(s interface{}, name string) any {
	f := reflect.ValueOf(s).Elem().FieldByName(name)
	return reflect.ValueOf(f).Interface()
}

func GetFieldValueAnCast[T any](s interface{}, name string) T {
	f := reflect.ValueOf(s).Elem().FieldByName(name)
	return reflect.ValueOf(f).Interface().(T)
}

func StructToMap(s any) map[string]any {
	var mp map[string]any
	j, _ := json.Marshal(s)
	json.Unmarshal(j, &mp)
	return mp
}

func CopyMapToStruct(m map[string]any, v interface{}) {
	out, _ := json.Marshal(m)
	json.Unmarshal(out, v)
}

/*func GetTag(s interface{}, fieldName string, tagName string) string {
	field, ok := reflect.TypeOf(s).FieldByName(fieldName)
	fmt.Printf("%v\n", fieldName)
	if ok {
		return fmt.Sprintf("%v", field.Tag)
	}
	return ""
}*/
