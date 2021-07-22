package cmhp

import (
	"fmt"
	"testing"
)

func TestAbs(t *testing.T) {
	x := Request(HttpArgs{
		Url:    "https://yandex.ru",
		Method: "GET",
	})
	fmt.Println(x.Body)
}
