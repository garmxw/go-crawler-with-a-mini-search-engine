// Package tui provides a simple interactive prompt-based form.
// It uses only fmt + bufio for input — no raw terminal mode, no bubbletea,
// no stdin hijacking. This makes it work identically on every terminal
// including Windows Terminal Preview, cmd.exe, and editor terminals.
package tui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ── Palette (lipgloss for output colours only — never touches stdin mode) ─────

var (
	clrCyan    = lipgloss.Color("#00D4FF")
	clrMagenta = lipgloss.Color("#FF2D78")
	clrGreen   = lipgloss.Color("#00FF9F")
	clrOrange  = lipgloss.Color("#FF7A00")
	clrGray    = lipgloss.Color("#5A5A7A")
	clrDimGray = lipgloss.Color("#2E2E42")
	clrWhite   = lipgloss.Color("#E8E8F0")
	clrDark    = lipgloss.Color("#0D0D1A")
	clrRed     = lipgloss.Color("#FF2D78")
	clrYellow  = lipgloss.Color("#FFD600")
)

var (
	sLabel   = lipgloss.NewStyle().Foreground(clrCyan).Bold(true).Width(14)
	sDefault = lipgloss.NewStyle().Foreground(clrGray)
	sValue   = lipgloss.NewStyle().Foreground(clrWhite).Bold(true)
	sPrompt  = lipgloss.NewStyle().Foreground(clrCyan).Bold(true)
	sError   = lipgloss.NewStyle().Foreground(clrRed).Bold(true)
	sDivider = lipgloss.NewStyle().Foreground(clrGray)
	sSuccess = lipgloss.NewStyle().Foreground(clrGreen).Bold(true)
	sSection = lipgloss.NewStyle().Foreground(clrCyan).Bold(true)

	sBadgeLocal = lipgloss.NewStyle().Background(clrGreen).Foreground(clrDark).Bold(true).Padding(0, 1)
	sBadgeWeb   = lipgloss.NewStyle().Background(clrCyan).Foreground(clrDark).Bold(true).Padding(0, 1)
	sBadgeLive  = lipgloss.NewStyle().Background(clrOrange).Foreground(clrDark).Bold(true).Padding(0, 1)

	sConfirmBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(clrCyan).
			Padding(0, 2)
)

// ── Low-level helpers ─────────────────────────────────────────────────────────

var scanner = bufio.NewScanner(os.Stdin)

// readLine reads one line from stdin (blocks until user presses Enter).
func readLine() string {
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

// ask prints a styled prompt and returns the trimmed response.
// If the user presses Enter without typing, defaultVal is returned.
func ask(label, defaultVal string) string {
	def := ""
	if defaultVal != "" {
		def = "  " + sDefault.Render("["+defaultVal+"]")
	}
	fmt.Print(sLabel.Render(label) + sPrompt.Render(" > ") + def + " ")
	raw := readLine()
	if raw == "" {
		return defaultVal
	}
	return raw
}

// askInt asks for an integer, loops until valid, returns defaultVal on empty Enter.
func askInt(label string, defaultVal, min, max int) int {
	for {
		raw := ask(label, strconv.Itoa(defaultVal))
		if raw == strconv.Itoa(defaultVal) {
			return defaultVal
		}
		n, err := strconv.Atoi(raw)
		if err != nil {
			fmt.Println(sError.Render("  [ER] must be a number"))
			continue
		}
		if n < min || (max >= 0 && n > max) {
			if max >= 0 {
				fmt.Println(sError.Render(fmt.Sprintf("  [ER] must be between %d and %d", min, max)))
			} else {
				fmt.Println(sError.Render(fmt.Sprintf("  [ER] must be >= %d", min)))
			}
			continue
		}
		return n
	}
}

// askChoice asks the user to pick from a numbered list, returns the chosen value.
func askChoice(label string, options []string, defaultIdx int) string {
	fmt.Println(sSection.Render("  " + label))
	for i, o := range options {
		marker := "  "
		if i == defaultIdx {
			marker = sSuccess.Render("* ")
		}
		fmt.Printf("  %s%d. %s\n", marker, i+1, sValue.Render(o))
	}
	for {
		raw := ask("choice", strconv.Itoa(defaultIdx+1))
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > len(options) {
			fmt.Println(sError.Render(fmt.Sprintf("  [ER] enter 1-%d", len(options))))
			continue
		}
		return options[n-1]
	}
}

// askBool asks yes/no, returns bool.
func askBool(label string, defaultVal bool) bool {
	def := "n"
	if defaultVal {
		def = "y"
	}
	for {
		raw := ask(label+" (y/n)", def)
		switch strings.ToLower(raw) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println(sError.Render("  [ER] enter y or n"))
		}
	}
}

// askURLs collects multiple URLs, one per line. Empty line stops input.
func askURLs(existingURLs []string) []string {
	urls := append([]string{}, existingURLs...)
	fmt.Println(sSection.Render("  URLs") + sDefault.Render("  (one per line, empty line to finish)"))
	for {
		raw := ask(fmt.Sprintf("url %d", len(urls)+1), "")
		if raw == "" {
			break
		}
		urls = append(urls, raw)
		fmt.Println(sSuccess.Render(fmt.Sprintf("    added: %s", raw)))
	}
	return urls
}

// divider prints a horizontal rule.
func divider() {
	fmt.Println(sDivider.Render(strings.Repeat("-", 56)))
}

// ── Config types ──────────────────────────────────────────────────────────────

type SearchConfig struct {
	Mode      string
	Query     string
	Path      string
	URLs      []string
	Lang      string
	Limit     int
	Detailed  bool
	File      string
	JSON      string
	Depth     int
	MaxPages  int
	Delay     int
	Storage   string
	Submitted bool
}

type CrawlerConfig struct {
	URLs      []string
	File      string
	JSON      string
	Depth     int
	MaxPages  int
	Delay     int
	Storage   string
	Submitted bool
}

// ── Search form ───────────────────────────────────────────────────────────────

func RunSearchForm(d SearchConfig) (SearchConfig, error) {
	if d.Limit <= 0 {
		d.Limit = 5
	}
	if d.MaxPages <= 0 {
		d.MaxPages = 3
	}
	if d.Delay <= 0 {
		d.Delay = 2
	}
	if d.Storage == "" {
		d.Storage = "data/pages"
	}
	if d.Mode == "" {
		d.Mode = "local"
	}
	if d.Lang == "" {
		d.Lang = "english"
	}

	fmt.Println()
	divider()
	fmt.Println(sSection.Render("  Search Configuration"))
	divider()
	fmt.Println()

	// Mode
	d.Mode = askChoice("Mode", []string{"local", "web", "live"}, modeIdx(d.Mode))
	fmt.Println()

	// Query
	for {
		d.Query = ask("Query", d.Query)
		if strings.TrimSpace(d.Query) != "" {
			break
		}
		fmt.Println(sError.Render("  [ER] query cannot be empty"))
	}
	fmt.Println()

	// Language
	d.Lang = askChoice("Language", []string{"english", "french"}, langIdx(d.Lang))
	fmt.Println()

	// Limit
	d.Limit = askInt("Limit", d.Limit, 1, 100)
	fmt.Println()

	// Detailed
	d.Detailed = askBool("Detailed output", d.Detailed)
	fmt.Println()

	// Mode-specific fields
	switch d.Mode {
	case "local":
		for {
			d.Path = ask("Local path", d.Path)
			if strings.TrimSpace(d.Path) != "" {
				break
			}
			fmt.Println(sError.Render("  [ER] path is required for local mode"))
		}
		fmt.Println()

	case "web":
		d.Storage = ask("Storage path", d.Storage)
		fmt.Println()

	case "live":
		d.URLs = askURLs(d.URLs)
		fmt.Println()
		d.File = ask("URL file", d.File)
		fmt.Println()
		d.JSON = ask("JSON file", d.JSON)
		fmt.Println()

		if len(d.URLs) == 0 && d.File == "" && d.JSON == "" {
			fmt.Println(sError.Render("  [ER] provide at least one URL, file, or JSON source"))
			d.URLs = askURLs(d.URLs)
			fmt.Println()
		}

		d.Depth = askInt("Depth", d.Depth, 0, 10)
		fmt.Println()
		d.MaxPages = askInt("Max pages", d.MaxPages, 1, -1)
		fmt.Println()
		d.Delay = askInt("Delay (s)", d.Delay, 0, 60)
		fmt.Println()
		d.Storage = ask("Storage path", d.Storage)
		fmt.Println()
	}

	// Confirm
	divider()
	printSearchConfirm(d)
	divider()
	fmt.Println()

	confirmed := askBool("Run with this config?", true)
	fmt.Println()

	if !confirmed {
		// Let user re-run from scratch by recursing once
		return RunSearchForm(d)
	}

	d.Submitted = true
	return d, nil
}

func printSearchConfirm(d SearchConfig) {
	var sb strings.Builder

	key := lipgloss.NewStyle().Foreground(clrGray).Width(12).Align(lipgloss.Right)
	val := lipgloss.NewStyle().Foreground(clrWhite).Bold(true)
	sep := lipgloss.NewStyle().Foreground(clrGray).Render("  |  ")

	row := func(k, v string) {
		sb.WriteString(key.Render(k) + sep + val.Render(v) + "\n")
	}

	var badge string
	switch d.Mode {
	case "local":
		badge = sBadgeLocal.Render(" LOCAL ")
	case "web":
		badge = sBadgeWeb.Render(" WEB ")
	case "live":
		badge = sBadgeLive.Render(" LIVE ")
	}

	sb.WriteString(sSection.Render("[*] Ready to run") + "\n")
	sb.WriteString(sDivider.Render(strings.Repeat("-", 44)) + "\n")
	sb.WriteString(key.Render("Mode") + sep + badge + "\n")
	row("Query", fmt.Sprintf("%q", d.Query))
	row("Language", d.Lang)
	row("Limit", strconv.Itoa(d.Limit))
	row("Detailed", strconv.FormatBool(d.Detailed))

	switch d.Mode {
	case "local":
		row("Path", d.Path)
	case "web":
		row("Storage", d.Storage)
	case "live":
		for i, u := range d.URLs {
			if i == 0 {
				row("URLs", u)
			} else {
				row("", u)
			}
		}
		if len(d.URLs) == 0 {
			row("URLs", "(from file/json)")
		}
		row("Depth", strconv.Itoa(d.Depth))
		row("Max Pages", strconv.Itoa(d.MaxPages))
		row("Delay (s)", strconv.Itoa(d.Delay))
		row("Storage", d.Storage)
	}

	fmt.Println(sConfirmBox.Render(sb.String()))
}

// ── Crawler form ──────────────────────────────────────────────────────────────

func RunCrawlerForm(d CrawlerConfig) (CrawlerConfig, error) {
	if d.MaxPages <= 0 {
		d.MaxPages = 3
	}
	if d.Delay <= 0 {
		d.Delay = 2
	}
	if d.Storage == "" {
		d.Storage = "data/pages"
	}

	fmt.Println()
	divider()
	fmt.Println(sSection.Render("  Crawler Configuration"))
	divider()
	fmt.Println()

	d.URLs = askURLs(d.URLs)
	fmt.Println()
	d.File = ask("URL file", d.File)
	fmt.Println()
	d.JSON = ask("JSON file", d.JSON)
	fmt.Println()

	if len(d.URLs) == 0 && d.File == "" && d.JSON == "" {
		fmt.Println(sError.Render("  [ER] provide at least one URL, file, or JSON source"))
		d.URLs = askURLs(d.URLs)
		fmt.Println()
	}

	d.Depth = askInt("Depth", d.Depth, 0, 10)
	fmt.Println()
	d.MaxPages = askInt("Max pages", d.MaxPages, 1, -1)
	fmt.Println()
	d.Delay = askInt("Delay (s)", d.Delay, 0, 60)
	fmt.Println()
	d.Storage = ask("Storage path", d.Storage)
	fmt.Println()

	// Confirm
	divider()
	printCrawlerConfirm(d)
	divider()
	fmt.Println()

	confirmed := askBool("Start crawl with this config?", true)
	fmt.Println()

	if !confirmed {
		return RunCrawlerForm(d)
	}

	d.Submitted = true
	return d, nil
}

func printCrawlerConfirm(d CrawlerConfig) {
	var sb strings.Builder

	key := lipgloss.NewStyle().Foreground(clrGray).Width(12).Align(lipgloss.Right)
	val := lipgloss.NewStyle().Foreground(clrWhite).Bold(true)
	sep := lipgloss.NewStyle().Foreground(clrGray).Render("  |  ")

	row := func(k, v string) {
		sb.WriteString(key.Render(k) + sep + val.Render(v) + "\n")
	}

	sb.WriteString(sSection.Render("[*] Ready to crawl") + "\n")
	sb.WriteString(sDivider.Render(strings.Repeat("-", 44)) + "\n")

	for i, u := range d.URLs {
		if i == 0 {
			row("URLs", u)
		} else {
			row("", u)
		}
	}
	if len(d.URLs) == 0 {
		row("URLs", "(from file/json)")
	}
	row("File", orDash(d.File))
	row("JSON", orDash(d.JSON))
	row("Depth", strconv.Itoa(d.Depth))
	row("Max Pages", strconv.Itoa(d.MaxPages))
	row("Delay (s)", strconv.Itoa(d.Delay))
	row("Storage", d.Storage)

	fmt.Println(sConfirmBox.Render(sb.String()))
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func modeIdx(mode string) int {
	for i, m := range []string{"local", "web", "live"} {
		if m == mode {
			return i
		}
	}
	return 0
}

func langIdx(lang string) int {
	if lang == "french" {
		return 1
	}
	return 0
}

func orDash(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}
