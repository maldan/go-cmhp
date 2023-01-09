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

func TestGetRange(t *testing.T) {
	newArray := []int{1, 2, 3, 4, 5}
	finalArray := cmhp_slice.GetRange(newArray, 1, 2)
	if len(finalArray) != 2 {
		t.Error("Fuck range")
	}

	finalArray = cmhp_slice.GetRange(newArray, 1, 10)
	if len(finalArray) != 4 {
		t.Error("Fuck range")
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

func TestPaginate(t *testing.T) {
	newArray := make([]int, 10)
	for i := 0; i < 10; i++ {
		newArray[i] = i
	}
	s := cmhp_slice.Paginate(newArray, 0, 10)
	if fmt.Sprintf("%v", s) != "[0 1 2 3 4 5 6 7 8 9]" {
		t.Error("Fuck paginate")
	}

	s = cmhp_slice.Paginate(newArray, 5, 10)
	if fmt.Sprintf("%v", s) != "[5 6 7 8 9]" {
		t.Error("Fuck paginate")
	}

	s = cmhp_slice.Paginate(newArray, 9, 10)
	if fmt.Sprintf("%v", s) != "[9]" {
		t.Error("Fuck paginate")
	}

	s = cmhp_slice.Paginate(newArray, 111, 111)
	if fmt.Sprintf("%v", s) != "[]" {
		t.Error("Fuck paginate")
	}

	s = cmhp_slice.Paginate(newArray, -1, 10)
	if fmt.Sprintf("%v", s) != "[0 1 2 3 4 5 6 7 8 9]" {
		t.Error("Fuck paginate")
	}

	s = cmhp_slice.Paginate(newArray, -1, -1)
	if fmt.Sprintf("%v", s) != "[]" {
		t.Error("Fuck paginate")
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

/*func BenchmarkTwo(b *testing.B) {
	newArray := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		cmhp_slice.FilterR(newArray, func(t interface{}) bool {
			return t.(int) > 3
		})
	}
}*/
