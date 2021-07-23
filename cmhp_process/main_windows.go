package cmhp_process

import (
	"os/exec"
	"strings"
	"syscall"
)

func Exec(args ...string) string {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c.Run()
	return b.String()
}

func Create(args ...string) (*exec.Cmd, *strings.Builder) {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c, b
}