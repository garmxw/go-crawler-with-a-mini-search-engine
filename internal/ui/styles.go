// Package ui provides all terminal UI utilities for the crawler and search CLI.
// Borders use NormalBorder/DoubleBorder only — RoundedBorder and ThickBorder
// use characters (╭╮╰╯ / ┃) that crash the Windows legacy console host.
package ui

import "github.com/charmbracelet/lipgloss"

// Colour palette

var (
	colorCyan    = lipgloss.Color("#00D4FF")
	colorMagenta = lipgloss.Color("#FF2D78")
	colorGreen   = lipgloss.Color("#00FF9F")
	colorYellow  = lipgloss.Color("#FFD600")
	colorOrange  = lipgloss.Color("#FF7A00")
	colorGold    = lipgloss.Color("#FFD700")
	colorSilver  = lipgloss.Color("#C0C0C0")
	colorBronze  = lipgloss.Color("#CD7F32")
	colorGray    = lipgloss.Color("#5A5A7A")
	colorDimGray = lipgloss.Color("#2E2E42")
	colorWhite   = lipgloss.Color("#E8E8F0")
	colorDark    = lipgloss.Color("#0D0D1A")
)

// Text styles

var (
	StyleTitle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(colorMagenta).
			Italic(true)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	StyleWarn = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(colorMagenta).
			Bold(true)

	StyleInfo = lipgloss.NewStyle().
			Foreground(colorCyan)

	StyleDim = lipgloss.NewStyle().
			Foreground(colorGray)

	StyleHighlight = lipgloss.NewStyle().
			Foreground(colorOrange).
			Bold(true)

	StyleLabel = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	StyleValue = lipgloss.NewStyle().
			Foreground(colorCyan)
)

// Box styles
// NormalBorder  (┌─┐│└┘) — safe on all Windows console hosts.
// DoubleBorder  (╔═╗║╚╝) — safe on all Windows console hosts.
// RoundedBorder (╭─╮│╰╯) — NOT safe; omitted.
// ThickBorder   (┏━┓┃┗┛) — NOT safe; replaced with DoubleBorder.

var (
	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorCyan).
			Padding(0, 2)

	StyleWarnBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorYellow).
			Padding(0, 2)

	StyleErrorBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorMagenta).
			Padding(0, 2)

	StyleSuccessBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorGreen).
			Padding(0, 2)

	StyleConfigBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorDimGray).
			Padding(0, 2).
			MarginTop(1)
)
