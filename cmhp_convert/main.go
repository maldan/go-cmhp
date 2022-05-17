package cmhp_convert

import (
	b64 "encoding/base64"
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

func ToBase64(v interface{}) string {
	switch v.(type) {
	case string:
		enc := b64.URLEncoding.EncodeToString([]byte(v.(string)))
		return enc
	default:
		return ""
	}
	return ""
}

func FromBase64(v string) []byte {
	uDec, _ := b64.URLEncoding.DecodeString(v)
	return uDec
}
