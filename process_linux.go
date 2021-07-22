package cmhp

import (
	"os/exec"
	"strings"
	"syscall"
)

func ProcessExec(args ...string) string {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.SysProcAttr = &syscall.SysProcAttr{
		Noctty: true,
	}
	c.Run()
	return b.String()
}

func ProcessCreate(args ...string) (*exec.Cmd, *strings.Builder) {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.SysProcAttr = &syscall.SysProcAttr{
		Noctty: true,
	}
	return c, b
}
