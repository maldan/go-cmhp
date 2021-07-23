package cmhp_process

import (
	"os/exec"
	"strings"
)

func Exec(args ...string) string {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.Run()
	c.Process.Release()
	return b.String()
}

func Create(args ...string) (*exec.Cmd, *strings.Builder) {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	return c, b
}
