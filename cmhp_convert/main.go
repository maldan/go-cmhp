package cmhp_convert

import (
	b64 "encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
	"sync"
)

func StrToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func StrToFloat(s string) float64 {
	n, _ := strconv.ParseFloat(s, 64)
	return n
}

func JsonToStruct[T any](s string) T {
	m := new(T)
	json.Unmarshal([]byte(s), &m)
	return *m
}

func StructToMap[T any](v *T) map[string]any {
	m := map[string]any{}
	bytes, _ := json.Marshal(v)
	json.Unmarshal(bytes, &m)
	return m
}

func StructToMapFast[T any](v *T) map[string]any {
	m := map[string]any{}

	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()
	for i := 0; i < typeOf.NumField(); i++ {
		m[typeOf.Field(i).Name] = valueOf.Field(i).Interface()
	}
	// fmt.Printf("%v\n", m)
	/*bytes, _ := json.Marshal(v)
	json.Unmarshal(bytes, &m)*/
	return m
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func MapToSyncMap(m map[interface{}]interface{}) sync.Map {
	sm := sync.Map{}
	for k, v := range m {
		sm.Store(k, v)
	}
	return sm
}

func SyncMapToMap(sm sync.Map) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	sm.Range(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

func ToUrlBase64[T string | []byte](v T) string {
	switch any(v).(type) {
	case string:
		enc := b64.URLEncoding.EncodeToString([]byte(any(v).(string)))
		return enc
	case []byte:
		enc := b64.URLEncoding.EncodeToString(any(v).([]byte))
		return enc
	default:
		return ""
	}
}

func ToStdBase64[T string | []byte](v T) string {
	switch any(v).(type) {
	case string:
		enc := b64.StdEncoding.EncodeToString([]byte(any(v).(string)))
		return enc
	case []byte:
		enc := b64.StdEncoding.EncodeToString(any(v).([]byte))
		return enc
	default:
		return ""
	}
}

func FromUrlBase64(v string) []byte {
	uDec, _ := b64.URLEncoding.DecodeString(v)
	return uDec
}

func FromStdBase64(v string) []byte {
	uDec, _ := b64.StdEncoding.DecodeString(v)
	return uDec
}
