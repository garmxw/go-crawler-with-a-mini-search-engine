package tui

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Palette ───────────────────────────────────────────────────────────────────

var (
	clrCyan    = lipgloss.Color("#00D4FF")
	clrMagenta = lipgloss.Color("#FF2D78")
	clrGreen   = lipgloss.Color("#00FF9F")
	clrYellow  = lipgloss.Color("#FFD600")
	clrOrange  = lipgloss.Color("#FF7A00")
	clrGray    = lipgloss.Color("#5A5A7A")
	clrDimGray = lipgloss.Color("#2E2E42")
	clrWhite   = lipgloss.Color("#E8E8F0")
	clrDark    = lipgloss.Color("#0D0D1A")
	clrRed     = lipgloss.Color("#FF2D78")
)

const formWidth = 64

// ── Borders ───────────────────────────────────────────────────────────────────
// NormalBorder  uses ┌┐└┘─│  — supported by every Windows console host.
// DoubleBorder  uses ╔╗╚╝═║  — supported by every Windows console host.
// RoundedBorder uses ╭╮╰╯    — NOT supported by the legacy conhost fallback;
// avoid it so Windows Terminal Preview doesn't crash.

var (
	// ── Form box ─────────────────────────────────────────────────────────────

	styleFormBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).       // ┌─┐  safe on Windows
			BorderForeground(clrCyan).
			Padding(1, 3).
			Width(formWidth)

	styleFormBoxLocked = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(clrGray).
				Padding(1, 3).
				Width(formWidth)

	// ── Typography ────────────────────────────────────────────────────────────

	styleFormTitle = lipgloss.NewStyle().
			Foreground(clrCyan).
			Bold(true)

	styleFormTitleLocked = lipgloss.NewStyle().
				Foreground(clrGray).
				Bold(true)

	styleFocusedLabel = lipgloss.NewStyle().
				Foreground(clrCyan).
				Bold(true).
				Width(16)

	styleUnfocusedLabel = lipgloss.NewStyle().
				Foreground(clrGray).
				Width(16)

	styleFocusedInput = lipgloss.NewStyle().
				Foreground(clrWhite).
				Background(clrDimGray).
				Padding(0, 1).
				Width(38)

	styleUnfocusedInput = lipgloss.NewStyle().
				Foreground(clrGray).
				Background(lipgloss.Color("#1A1A2E")).
				Padding(0, 1).
				Width(38)

	styleCursor = lipgloss.NewStyle().
			Background(clrCyan).
			Foreground(clrDark).
			Bold(true)

	styleDim = lipgloss.NewStyle().
			Foreground(clrGray)

	styleHint = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3A3A5A")).
			Italic(true)

	styleError = lipgloss.NewStyle().
			Foreground(clrRed).
			Bold(true)

	// ── Badges ────────────────────────────────────────────────────────────────

	styleSelectedBadge = lipgloss.NewStyle().
				Background(clrCyan).
				Foreground(clrDark).
				Bold(true).
				Padding(0, 1)

	styleUnselectedBadge = lipgloss.NewStyle().
				Background(clrDimGray).
				Foreground(clrGray).
				Padding(0, 1)

	styleOffBadge = lipgloss.NewStyle().
			Background(clrMagenta).
			Foreground(clrDark).
			Bold(true).
			Padding(0, 1)

	styleOnBadge = lipgloss.NewStyle().
			Background(clrGreen).
			Foreground(clrDark).
			Bold(true).
			Padding(0, 1)

	styleURLTag = lipgloss.NewStyle().
			Background(lipgloss.Color("#1A2A3A")).
			Foreground(clrCyan).
			Padding(0, 1)

	// ── Buttons ───────────────────────────────────────────────────────────────

	styleSubmitActive = lipgloss.NewStyle().
				Background(clrGreen).
				Foreground(clrDark).
				Bold(true).
				Padding(0, 4)

	styleSubmitInactive = lipgloss.NewStyle().
				Background(clrDimGray).
				Foreground(clrGray).
				Padding(0, 4)

	styleSubmitDone = lipgloss.NewStyle().
			Background(clrGray).
			Foreground(clrDark).
			Bold(true).
			Padding(0, 4)

	// ── Confirm panel ─────────────────────────────────────────────────────────

	styleConfirmBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).       // ╔═╗  safe on Windows
			BorderForeground(clrCyan).
			Padding(0, 2)

	styleConfirmTitle = lipgloss.NewStyle().
				Foreground(clrCyan).
				Bold(true)
)

// ── Key hint bar ──────────────────────────────────────────────────────────────

func renderKeyHints(pairs ...string) string {
	keyStyle := lipgloss.NewStyle().
		Background(clrDimGray).
		Foreground(clrCyan).
		Bold(true).
		Padding(0, 1)
	descStyle := lipgloss.NewStyle().Foreground(clrGray)

	var parts []string
	for i := 0; i+1 < len(pairs); i += 2 {
		parts = append(parts, keyStyle.Render(pairs[i])+" "+descStyle.Render(pairs[i+1]))
	}
	return strings.Join(parts, "  ")
}

// ── Divider ───────────────────────────────────────────────────────────────────

func divider() string {
	// Use plain hyphen-minus instead of Unicode ─ (U+2500) for Windows safety.
	return styleDim.Render(strings.Repeat("-", formWidth-8))
}

// ── Program constructor ───────────────────────────────────────────────────────

// newProgram creates a Bubbletea program compatible with Windows Terminal Preview
// and all Unix terminals.
//
// tea.WithInput(os.Stdin): tells bubbletea to use the process stdin directly
// instead of opening /dev/tty, which does not exist on Windows and causes an
// immediate panic in Windows Terminal Preview.
//
// We do NOT use tea.WithOutput(os.Stderr): bubbletea putting stderr into raw
// mode corrupts the slog output that the crawler also writes to stderr, and on
// Windows Terminal Preview this causes colly's HTTP goroutines to fail silently,
// resulting in 0 pages crawled. TUI writes to stdout (default) and slog writes
// to stderr (default) — they never share a file descriptor.
//
// No WithAltScreen — the form stays visible in the normal scroll buffer.
func newProgram(m tea.Model) *tea.Program {
	return tea.NewProgram(
		m,
		tea.WithInput(os.Stdin),
	)
}
