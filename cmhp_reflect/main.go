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
