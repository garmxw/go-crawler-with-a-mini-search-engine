//go:build windows

package ui

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	for _, f := range []*os.File{os.Stdout, os.Stderr} {
		handle := windows.Handle(f.Fd())
		var mode uint32
		if err := windows.GetConsoleMode(handle, &mode); err != nil {
			continue
		}
		_ = windows.SetConsoleMode(handle, mode|0x0004|0x0001)
	}
}
