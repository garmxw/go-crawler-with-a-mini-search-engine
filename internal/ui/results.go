package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

// SearchResult mirrors models.SearchResult — keeps ui dependency-free.
type SearchResult struct {
	DocID int
	Path  string
	Score float64
}

// PrintSearchResults renders a ranked, scored results list.
// Bar uses '#' and '.' — safe on every Windows console host.
// Block chars like █ and ░ are NOT in CP437 on all Windows versions.
func PrintSearchResults(results []SearchResult, query string) {
	if len(results) == 0 {
		fmt.Println()
		fmt.Println(StyleWarnBox.Render(
			StyleWarn.Render("[!!]  No results found for  ") +
				StyleHighlight.Render("\""+query+"\""),
		))
		fmt.Println()
		return
	}

	SectionHeader(fmt.Sprintf("Results for \"%s\" (%d found)", query, len(results)))

	maxScore := results[0].Score
	barWidth := 22

	scoreStyle := lipgloss.NewStyle().
		Foreground(colorYellow).
		Bold(true).
		Width(8)

	pathStyle := lipgloss.NewStyle().Foreground(colorWhite)
	dimStyle := StyleDim
	idStyle := lipgloss.NewStyle().Foreground(colorGray)

	for i, r := range results {
		var rankBg lipgloss.Color
		switch i {
		case 0:
			rankBg = colorGold
		case 1:
			rankBg = colorSilver
		case 2:
			rankBg = colorBronze
		default:
			rankBg = colorDimGray
		}

		rankBadge := lipgloss.NewStyle().
			Foreground(colorDark).
			Background(rankBg).
			Bold(true).
			Padding(0, 1).
			Width(4).
			Align(lipgloss.Center)

		normalised := r.Score / maxScore
		filled := int(normalised * float64(barWidth))
		if filled < 1 {
			filled = 1
		}
		if filled > barWidth {
			filled = barWidth
		}

		// '#' for filled, '.' for empty — safe on all Windows consoles
		barFilled := lipgloss.NewStyle().Foreground(colorGreen).Render(strings.Repeat("#", filled))
		barEmpty := dimStyle.Render(strings.Repeat(".", barWidth-filled))
		bar := "[" + barFilled + barEmpty + "]"

		// Trim long paths — use "..." instead of ellipsis char
		displayPath := r.Path
		if len(displayPath) > 55 {
			displayPath = "..." + displayPath[len(displayPath)-52:]
		}

		fmt.Printf(
			"  %s  %s  %s  %s\n",
			rankBadge.Render(fmt.Sprintf("#%d", i+1)),
			scoreStyle.Render(fmt.Sprintf("%.4f", r.Score)),
			bar,
			pathStyle.Render(displayPath),
		)
		fmt.Println(idStyle.Render(fmt.Sprintf("             doc id: %d", r.DocID)))
		fmt.Println()
	}
}

// PrintCrawlSummary renders a summary table after a crawl finishes.
func PrintCrawlSummary(totalPages int, storagePath string) {
	SectionHeader("Crawl Summary")

	data := pterm.TableData{
		{pterm.FgCyan.Sprint("Metric"), pterm.FgCyan.Sprint("Value")},
		{"Pages crawled", pterm.FgGreen.Sprintf("%d", totalPages)},
		{"Storage path", pterm.FgYellow.Sprint(storagePath)},
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed(true).
		WithData(data).
		Render()

	fmt.Println()
}

// PrintIndexSummary renders a summary table after indexing finishes.
func PrintIndexSummary(totalDocs int, uniqueTokens int) {
	SectionHeader("Index Summary")

	data := pterm.TableData{
		{pterm.FgCyan.Sprint("Metric"), pterm.FgCyan.Sprint("Value")},
		{"Documents indexed", pterm.FgGreen.Sprintf("%d", totalDocs)},
		{"Unique tokens", pterm.FgYellow.Sprintf("%d", uniqueTokens)},
	}

	pterm.DefaultTable.
		WithHasHeader().
		WithBoxed(true).
		WithData(data).
		Render()

	fmt.Println()
}
