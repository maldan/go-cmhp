package cmhp_convert

import (
	"strconv"
	"sync"
)

func StrToInt(s string) int {
	n, _ := strconv.Atoi(s)
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
