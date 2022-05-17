package cmhp_slice

import (
	"math/rand"
	"time"
)

func Includes[T comparable](slice []T, v T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			return true
		}
	}

	return false
}

func Unique[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	list := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if _, value := keys[slice[i]]; !value {
			keys[slice[i]] = true
			list = append(list, slice[i])
		}
	}

	return list
}

func Filter[T comparable](slice []T, filter func(T) bool) []T {
	filtered := make([]T, 0)

	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

//func Sort[T comparable](slice []T, s func(i, j T) int) []T {
//sorted := make([]T, 0)
//return sorted
//copy := Filter(slice, func(i T) bool { return true })
//sort.SliceStable(copy, s)
//return copy
//panic("xxx")
//return slice
//}
//

func Map[T any, R any](slice []T, m func(T) R) []R {
	mapped := make([]R, 0)
	for _, v := range slice {
		mapped = append(mapped, m(v))
	}
	return mapped
}

var r = rand.New(rand.NewSource(time.Now().Unix()))

func PickRandom[T any](slice []T) T {
	return slice[r.Intn(len(slice))]
}
