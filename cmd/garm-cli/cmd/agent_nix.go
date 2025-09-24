//go:build !windows
// +build !windows

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/term"
)

var sigs = make(chan os.Signal, 1)

func watchTermResize(resizeCh chan [2]int, sessionID uuid.UUID) {
	signal.Notify(sigs, syscall.SIGWINCH)

	for range sigs {
		w, h, err := term.GetSize(int(os.Stdin.Fd()))
		if err == nil && sessionID != uuid.Nil {
			resizeCh <- [2]int{w, h}
		}
	}
}
