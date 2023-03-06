//go:build !windows
// +build !windows

package cmd

import (
	"os"
	"syscall"
)

var signals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
}
