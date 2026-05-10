//go:build windows

package tui

import (
	"os"

	"golang.org/x/sys/windows"
)

// originalStdinMode holds the console input mode before the TUI takes over.
// We restore it after p.Run() so that subsequent code (the crawler, slog, colly)
// gets a normal blocking console — not bubbletea's raw non-blocking mode.
var originalStdinMode uint32
var originalStdinSaved bool

func init() {
	// 1. Save stdin's current mode so we can restore it after the TUI exits.
	stdinHandle := windows.Handle(os.Stdin.Fd())
	if err := windows.GetConsoleMode(stdinHandle, &originalStdinMode); err == nil {
		originalStdinSaved = true
	}

	// 2. Enable VTP on stdout and stderr so lipgloss ANSI codes render correctly.
	for _, f := range []*os.File{os.Stdout, os.Stderr} {
		handle := windows.Handle(f.Fd())
		var mode uint32
		if err := windows.GetConsoleMode(handle, &mode); err != nil {
			continue
		}
		// ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
		// ENABLE_PROCESSED_OUTPUT            = 0x0001
		_ = windows.SetConsoleMode(handle, mode|0x0004|0x0001)
	}
}

// restoreConsoleInput puts stdin back into the mode it was in before the TUI.
// Bubbletea sets stdin to raw mode (no echo, no line buffering, non-blocking).
// On Windows Terminal Preview this prevents colly's HTTP goroutines from
// blocking on network I/O correctly, so every request times out instantly and
// total_pages=0. Restoring the original mode fixes this.
func restoreConsoleInput() {
	if !originalStdinSaved {
		return
	}
	stdinHandle := windows.Handle(os.Stdin.Fd())
	_ = windows.SetConsoleMode(stdinHandle, originalStdinMode)
}
