package cmhp

import (
	"fmt"
	"testing"

	"github.com/maldan/go-cmhp/cmhp_net"
)

func TestAbs(t *testing.T) {
	x := cmhp_net.Request(cmhp_net.HttpArgs{
		Url:    "https://yandex.ru",
		Method: "GET",
	})
	fmt.Println(x.Body)
}
