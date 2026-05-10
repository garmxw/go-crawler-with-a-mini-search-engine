//go:build !windows

package tui

// restoreConsoleInput is a no-op on non-Windows platforms.
// On Windows it restores stdin's console mode after the TUI exits.
func restoreConsoleInput() {}
