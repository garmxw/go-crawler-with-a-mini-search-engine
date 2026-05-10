package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Inline messages
// Symbols use only characters in Windows-1252 / CP437 that every console host
// can render. Avoided: ✔ ✘ ⚠ ◆ ▸ ⚙ │ — ─ (replaced with ASCII equivalents).

func Success(msg string) {
	fmt.Println(StyleSuccess.Render("  [OK]  " + msg))
}

func Warn(msg string) {
	fmt.Println(StyleWarn.Render("  [!!]  " + msg))
}

func Error(msg string) {
	fmt.Println(StyleError.Render("  [ER]  " + msg))
}

func Info(msg string) {
	fmt.Println(StyleInfo.Render("  [>>]  " + msg))
}

func Dim(msg string) {
	fmt.Println(StyleDim.Render("        " + msg))
}

// Boxed messages

func SuccessBox(msg string) {
	fmt.Println()
	fmt.Println(StyleSuccessBox.Render(StyleSuccess.Render("[OK]  " + msg)))
	fmt.Println()
}

func WarnBox(msg string) {
	fmt.Println()
	fmt.Println(StyleWarnBox.Render(StyleWarn.Render("[!!]  " + msg)))
	fmt.Println()
}

func ErrorBox(msg string) {
	fmt.Println()
	fmt.Println(StyleErrorBox.Render(StyleError.Render("[ER]  " + msg)))
	fmt.Println()
}

func InfoBox(msg string) {
	fmt.Println()
	fmt.Println(StyleBox.Render(StyleInfo.Render("[>>]  " + msg)))
	fmt.Println()
}

// Section header

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
	fmt.Println(style.Render("  >>  " + title))
	fmt.Println()
}

// Config panel

type ConfigRow struct {
	Key   string
	Value string
}

func PrintConfigPanel(title string, rows []ConfigRow) {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(colorCyan).
		Bold(true)

	sb.WriteString(titleStyle.Render("[*] "+title) + "\n")
	sb.WriteString(StyleDim.Render(strings.Repeat("-", 42)) + "\n")

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
				sepStyle.Render("  |  ") +
				valStyle.Render(row.Value) +
				"\n",
		)
	}

	fmt.Println(StyleConfigBox.Render(sb.String()))
}

// Mode badge

func PrintModeBadge(mode string) {
	var bg lipgloss.Color
	var label string

	switch mode {
	case "local":
		bg = colorGreen
		label = "[LOCAL]"
	case "web":
		bg = colorCyan
		label = "[WEB]"
	case "live":
		bg = colorOrange
		label = "[LIVE]"
	default:
		bg = colorGray
		label = "[" + strings.ToUpper(mode) + "]"
	}

	badge := lipgloss.NewStyle().
		Background(bg).
		Foreground(colorDark).
		Bold(true).
		Padding(0, 2)

	prefix := lipgloss.NewStyle().Foreground(colorGray)

	fmt.Println()
	fmt.Println(prefix.Render("  Mode  ") + badge.Render(label))
	fmt.Println()
}

// Single stat row

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

// Divider

func Divider() {
	fmt.Println(StyleDim.Render(strings.Repeat("-", 64)))
}
