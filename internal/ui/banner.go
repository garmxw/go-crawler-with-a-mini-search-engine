package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const crawlerASCII = `
                                          .__
   ____   ____   ________________ __  _  _|  |
  / ___\ /  _ \_/ ___\_  __ \__  \\ \/ \/ /  |
 / /_/  >  <_> )  \___|  | \// __ \\     /|  |__
 \___  / \____/ \___  >__|  (____  /\/\_/ |____/
/_____/             \/           \/             `

const searchASCII = `
                                                 .__
   ____   ____  ______ ____ _____ _______   ____ │  │__
  ╱ ___╲ ╱  _ ╲╱  ___╱╱ __ ╲╲__  ╲╲_  __ ╲_╱ ___╲│  │  ╲
 ╱ ╱_╱  >  <_> )___ ╲╲  ___╱ ╱ __ ╲│  │ ╲╱╲  ╲___│   Y  ╲
 ╲___  ╱ ╲____╱____  >╲___  >____  ╱__│    ╲___  >___│  ╱
╱_____╱            ╲╱     ╲╱     ╲╱            ╲╱     ╲╱`

// PrintCrawlerBanner prints the styled crawler banner.
func PrintCrawlerBanner() {
	asciiStyle := lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	taglineStyle := lipgloss.NewStyle().Foreground(colorMagenta).Italic(true)
	versionBadge := lipgloss.NewStyle().
		Foreground(colorDark).
		Background(colorCyan).
		Bold(true).
		Padding(0, 1)

	fmt.Println(asciiStyle.Render(crawlerASCII))

	center := lipgloss.NewStyle().Width(64).Align(lipgloss.Center)
	fmt.Println(center.Render(
		taglineStyle.Render("High-speed web crawler") +
			"   " +
			versionBadge.Render(" v1.0 "),
	))
	fmt.Println(StyleDim.Render(strings.Repeat("-", 64)))
	fmt.Println()
}

// PrintSearchBanner prints the styled search engine banner.
func PrintSearchBanner() {
	asciiStyle := lipgloss.NewStyle().Foreground(colorMagenta).Bold(true)
	taglineStyle := lipgloss.NewStyle().Foreground(colorCyan).Italic(true)
	versionBadge := lipgloss.NewStyle().
		Foreground(colorDark).
		Background(colorMagenta).
		Bold(true).
		Padding(0, 1)

	fmt.Println(asciiStyle.Render(searchASCII))

	center := lipgloss.NewStyle().Width(64).Align(lipgloss.Center)
	fmt.Println(center.Render(
		taglineStyle.Render("TF-IDF powered mini search engine") +
			"   " +
			versionBadge.Render(" v1.0 "),
	))
	fmt.Println(StyleDim.Render(strings.Repeat("-", 64)))
	fmt.Println()
}
