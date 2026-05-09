package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Inline messages

// Success prints a green success message with a checkmark.
func Success(msg string) {
	fmt.Println(StyleSuccess.Render("  ✔  " + msg))
}

// Warn prints a yellow warning message.
func Warn(msg string) {
	fmt.Println(StyleWarn.Render("  ⚠  " + msg))
}

// Error prints a red error message.
func Error(msg string) {
	fmt.Println(StyleError.Render("  ✘  " + msg))
}

// Info prints a cyan informational message.
func Info(msg string) {
	fmt.Println(StyleInfo.Render("  ◆  " + msg))
}

// Dim prints a muted gray message.
func Dim(msg string) {
	fmt.Println(StyleDim.Render("     " + msg))
}

// Boxed messages

// SuccessBox renders a green-bordered success box.
func SuccessBox(msg string) {
	fmt.Println()
	fmt.Println(StyleSuccessBox.Render(StyleSuccess.Render("✔  " + msg)))
	fmt.Println()
}

// WarnBox renders a yellow-bordered warning box.
func WarnBox(msg string) {
	fmt.Println()
	fmt.Println(StyleWarnBox.Render(StyleWarn.Render("⚠  " + msg)))
	fmt.Println()
}

// ErrorBox renders a magenta-bordered error box.
func ErrorBox(msg string) {
	fmt.Println()
	fmt.Println(StyleErrorBox.Render(StyleError.Render("✘  " + msg)))
	fmt.Println()
}

// InfoBox renders a cyan-bordered info box.
func InfoBox(msg string) {
	fmt.Println()
	fmt.Println(StyleBox.Render(StyleInfo.Render("◆  " + msg)))
	fmt.Println()
}

// Section header

// SectionHeader prints a bold underlined section title.
func SectionHeader(title string) {
	style := lipgloss.NewStyle().
		Foreground(colorCyan).
		Bold(true).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorGray).
		Width(60).
		PaddingBottom(0)

	fmt.Println()
	fmt.Println(style.Render("  ▸  " + title))
	fmt.Println()
}

// Config panel

// ConfigRow is a single key–value row in the config panel.
type ConfigRow struct {
	Key   string
	Value string
}

// PrintConfigPanel renders a styled configuration panel with key–value rows.
func PrintConfigPanel(title string, rows []ConfigRow) {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(colorCyan).
		Bold(true)

	sb.WriteString(titleStyle.Render("⚙  "+title) + "\n")
	sb.WriteString(StyleDim.Render(strings.Repeat("─", 42)) + "\n")

	keyStyle := lipgloss.NewStyle().
		Foreground(colorGray).
		Width(14).
		Align(lipgloss.Right)

	sepStyle := StyleDim
	valStyle := lipgloss.NewStyle().
		Foreground(colorWhite).
		Bold(true)

	for _, row := range rows {
		sb.WriteString(
			keyStyle.Render(row.Key) +
				sepStyle.Render("  │  ") +
				valStyle.Render(row.Value) +
				"\n",
		)
	}

	fmt.Println(StyleConfigBox.Render(sb.String()))
}

// Mode badge

// PrintModeBadge prints a coloured pill-style badge indicating the active mode.
func PrintModeBadge(mode string) {
	var bg lipgloss.Color
	var icon string

	switch mode {
	case "local":
		bg = colorGreen
		icon = "📁"
	case "web":
		bg = colorCyan
		icon = "🌐"
	case "live":
		bg = colorOrange
		icon = "⚡"
	default:
		bg = colorGray
		icon = "•"
	}

	badge := lipgloss.NewStyle().
		Background(bg).
		Foreground(colorDark).
		Bold(true).
		Padding(0, 2)

	prefix := lipgloss.NewStyle().Foreground(colorGray)

	fmt.Println()
	fmt.Println(
		prefix.Render("  Mode  ") +
			badge.Render(icon+"  "+strings.ToUpper(mode)),
	)
	fmt.Println()
}

// Single stat row

// PrintStat prints a single labelled metric.
func PrintStat(label string, value interface{}) {
	lStyle := lipgloss.NewStyle().
		Foreground(colorGray).
		Width(22)

	vStyle := lipgloss.NewStyle().
		Foreground(colorCyan).
		Bold(true)

	var vStr string
	switch v := value.(type) {
	case int:
		vStr = strconv.Itoa(v)
	case float64:
		vStr = fmt.Sprintf("%.4f", v)
	case string:
		vStr = v
	default:
		vStr = fmt.Sprintf("%v", v)
	}

	fmt.Println("  " + lStyle.Render(label) + vStyle.Render(vStr))
}

//Divider

// Divider prints a horizontal rule.
func Divider() {
	fmt.Println(StyleDim.Render(strings.Repeat("─", 64)))
}
