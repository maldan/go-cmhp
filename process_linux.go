package cmhp

import (
	"os/exec"
	"strings"
)

func ProcessExec(args ...string) string {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.Run()
	return b.String()
}

func ProcessCreate(args ...string) (*exec.Cmd, *strings.Builder) {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	return c, b
}
