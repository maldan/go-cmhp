package cmhp_slice_test

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"testing"
)

func TestIncludes(t *testing.T) {
	newArray := []int{1, 2, 3, 4, 5}
	if !cmhp_slice.Includes(newArray, 2) {
		t.Error("Fuck includes")
	}
}

func TestUnique(t *testing.T) {
	newArray := []string{"x", "y", "z", "x", "xx", "d"}
	if len(cmhp_slice.Unique(newArray)) != 5 {
		t.Error("Fuck unique")
	}
}

func TestFilter(t *testing.T) {
	newArray := []int{1, 2, 3, 4, 5}
	finalArray := cmhp_slice.Filter(newArray, func(t int) bool {
		return t > 3
	})
	if len(finalArray) != 2 {
		t.Error("Fuck filter")
	}
}

func TestMap(t *testing.T) {
	newArray := []int{1, 2, 3, 4, 5}
	finalArray := cmhp_slice.Map(newArray, func(t int) string {
		return fmt.Sprintf("%v", t)
	})
	if finalArray[0] != "1" {
		t.Error("Fuck filter")
	}
}

func BenchmarkOne(b *testing.B) {
	newArray := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		cmhp_slice.Filter(newArray, func(t int) bool {
			return t > 3
		})
	}
}

func BenchmarkTwo(b *testing.B) {
	newArray := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		cmhp_slice.FilterR(newArray, func(t interface{}) bool {
			return t.(int) > 3
		})
	}
}
