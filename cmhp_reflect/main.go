package cmhp_reflect

import "reflect"

func SetField(s interface{}, name string, v interface{}) {
	f := reflect.ValueOf(s).Elem().FieldByName(name)
	if f.CanSet() {
		f.Set(reflect.ValueOf(v))
	}
}
