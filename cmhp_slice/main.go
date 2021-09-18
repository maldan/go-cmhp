package cmhp_slice

import (
	"reflect"
	"sort"
)

func Includes(slice []interface{}, v interface{}) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return true
		}
	}

	return false
}

func IncludesR(slice interface{}, v interface{}) bool {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == v {
			return true
		}
	}

	return false
}

func FindIndexR(slice interface{}, find func(interface{}) bool) int {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if find(s.Index(i).Interface()) {
			return i
		}
	}

	return -1
}

func FindR(slice interface{}, find func(interface{}) bool) (interface{}, int) {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if find(s.Index(i).Interface()) {
			return s.Index(i).Interface(), i
		}
	}

	return nil, -1
}

func FilterR(slice interface{}, filter func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0)

	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if filter(s.Index(i).Interface()) {
			filtered = append(filtered, s.Index(i).Interface())
		}
	}

	return filtered
}

func Filter(slice []interface{}, filter func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0)

	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Map(slice []interface{}, m func(interface{}) interface{}) []interface{} {
	mapped := make([]interface{}, 0)
	for _, v := range slice {
		mapped = append(mapped, m(v))
	}
	return mapped
}

func Sort(slice []interface{}, s func(i, j int) bool) []interface{} {
	copy := Filter(slice, func(i interface{}) bool { return true })
	sort.SliceStable(copy, s)
	return copy
}
