//go:build windows && !linux
// +build windows,!linux

package cmd

import "os"

var signals = []os.Signal{
	os.Interrupt,
}
