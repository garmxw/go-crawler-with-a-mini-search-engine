package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Palette

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

var (
	styleFormBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(clrCyan).
			Padding(1, 3).
			Width(formWidth)

	styleFormBoxLocked = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(clrGray).
				Padding(1, 3).
				Width(formWidth)

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

	styleConfirmBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(clrCyan).
			Padding(0, 2)

	styleConfirmTitle = lipgloss.NewStyle().
				Foreground(clrCyan).
				Bold(true)
)

// renderKeyHints builds a compact key binding bar from key/desc pairs.
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

func divider() string {
	return styleDim.Render(strings.Repeat("─", formWidth-8))
}
