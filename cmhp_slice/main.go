package cmhp_slice

import (
	"math/rand"
	"time"
)

func IndexOf[T any](slice []T, predicate func(T) bool) int {
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			return i
		}
	}

	return -1
}

func RemoveAt[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

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

func UniqueBy[T comparable](slice []T, fn func(T) any) []T {
	keys := make(map[any]T)
	list := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		keys[fn(slice[i])] = slice[i]
	}

	for _, v := range keys {
		list = append(list, v)
	}

	return list
}

func Filter[T any](slice []T, filter func(T) bool) []T {
	filtered := make([]T, 0)

	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func GetRange[T any](slice []T, fromIndex int, length int) []T {
	filtered := make([]T, 0)
	for i := 0; i < length; i++ {
		if fromIndex+i >= len(slice) {
			break
		}
		filtered = append(filtered, slice[fromIndex+i])
	}
	return filtered
}

func GetMapKeys[K comparable, V comparable](mmap map[K]V) []K {
	l := make([]K, 0)
	for k, _ := range mmap {
		l = append(l, k)
	}
	return l
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

func MaxSize[T any](slice []T, max int) []T {
	filtered := make([]T, 0)
	for i, v := range slice {
		if i >= max {
			break
		}
		filtered = append(filtered, v)
	}
	return filtered
}
