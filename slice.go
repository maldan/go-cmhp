package cmhp

import (
	"reflect"
	"sort"
)

func SliceIncludes(slice []interface{}, v interface{}) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return true
		}
	}

	return false
}

func SliceFindIndexR(slice interface{}, find func(interface{}) bool) int {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if find(s.Index(i).Interface()) {
			return i
		}
	}

	return -1
}

func SliceFindR(slice interface{}, find func(interface{}) bool) (interface{}, int) {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if find(s.Index(i).Interface()) {
			return s.Index(i).Interface(), i
		}
	}

	return nil, -1
}

func SliceFilterR(slice interface{}, filter func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0)

	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if filter(s.Index(i).Interface()) {
			filtered = append(filtered, s.Index(i).Interface())
		}
	}

	return filtered
}

func SliceFilter(slice []interface{}, filter func(interface{}) bool) []interface{} {
	filtered := make([]interface{}, 0)

	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func SliceMap(slice []interface{}, m func(interface{}) interface{}) []interface{} {
	mapped := make([]interface{}, 0)
	for _, v := range slice {
		mapped = append(mapped, m(v))
	}
	return mapped
}

func SliceSort(slice []interface{}, s func(i, j int) bool) []interface{} {
	copy := SliceFilter(slice, func(i interface{}) bool { return true })
	sort.SliceStable(copy, s)
	return copy
}
