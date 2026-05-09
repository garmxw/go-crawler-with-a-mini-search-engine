package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Config ────────────────────────────────────────────────────────────────────

type CrawlerConfig struct {
	URLs      []string
	File      string
	JSON      string
	Depth     int
	MaxPages  int
	Delay     int // crawl delay in seconds, default 2
	Storage   string
	Submitted bool
}

// ── Stage ─────────────────────────────────────────────────────────────────────

type crawlerStage int

const (
	cStageForm    crawlerStage = iota
	cStageConfirm
	cStageDone
)

// ── Field indices ─────────────────────────────────────────────────────────────

const (
	cfURLs = iota
	cfFile
	cfJSON
	cfDepth
	cfMaxPages
	cfDelay
	cfStorage
	cfSubmit
)

var crawlerFieldOrder = []int{cfURLs, cfFile, cfJSON, cfDepth, cfMaxPages, cfDelay, cfStorage, cfSubmit}

// ── Model ─────────────────────────────────────────────────────────────────────

type CrawlerFormModel struct {
	urlsInput    MultiInput
	fileInput    TextInput
	jsonInput    TextInput
	depthStep    Stepper
	maxPagesStep Stepper
	delayStep    Stepper
	storageInput TextInput

	focus  int
	stage  crawlerStage
	errMsg string
	quit   bool
}

func NewCrawlerForm(d CrawlerConfig) CrawlerFormModel {
	if d.MaxPages <= 0 {
		d.MaxPages = 3
	}
	if d.Delay <= 0 {
		d.Delay = 2
	}
	if d.Storage == "" {
		d.Storage = "data/pages"
	}

	m := CrawlerFormModel{
		urlsInput:    NewMultiInput("URLs", "https://example.com"),
		fileInput:    NewTextInput("URL File", "urls.txt"),
		jsonInput:    NewTextInput("JSON File", "urls.json"),
		depthStep:    NewStepper("Depth", d.Depth, 0, 10),
		maxPagesStep: NewStepper("Max Pages", d.MaxPages, 1, -1),
		delayStep:    NewStepper("Delay (s)", d.Delay, 0, 60),
		storageInput: NewTextInput("Storage", "data/pages"),
	}

	m.urlsInput.Values = d.URLs
	m.fileInput.Value = d.File
	m.fileInput.cursor = len([]rune(d.File))
	m.jsonInput.Value = d.JSON
	m.jsonInput.cursor = len([]rune(d.JSON))
	m.storageInput.Value = d.Storage
	m.storageInput.cursor = len([]rune(d.Storage))

	m.focusField(cfURLs)
	return m
}

func (m CrawlerFormModel) isTextFocused() bool {
	switch m.focus {
	case cfURLs, cfFile, cfJSON, cfStorage:
		return true
	}
	return false
}

func (m *CrawlerFormModel) focusField(idx int) {
	m.urlsInput.Blur()
	m.fileInput.Blur()
	m.jsonInput.Blur()
	m.depthStep.Blur()
	m.maxPagesStep.Blur()
	m.delayStep.Blur()
	m.storageInput.Blur()

	m.focus = idx
	switch idx {
	case cfURLs:
		m.urlsInput.Focus()
	case cfFile:
		m.fileInput.Focus()
	case cfJSON:
		m.jsonInput.Focus()
	case cfDepth:
		m.depthStep.Focus()
	case cfMaxPages:
		m.maxPagesStep.Focus()
	case cfDelay:
		m.delayStep.Focus()
	case cfStorage:
		m.storageInput.Focus()
	}
}

func (m CrawlerFormModel) nextField() int {
	for i, f := range crawlerFieldOrder {
		if f == m.focus && i+1 < len(crawlerFieldOrder) {
			return crawlerFieldOrder[i+1]
		}
	}
	return crawlerFieldOrder[0]
}

func (m CrawlerFormModel) prevField() int {
	for i, f := range crawlerFieldOrder {
		if f == m.focus && i > 0 {
			return crawlerFieldOrder[i-1]
		}
	}
	return crawlerFieldOrder[len(crawlerFieldOrder)-1]
}

// ── Bubbletea ─────────────────────────────────────────────────────────────────

func (m CrawlerFormModel) Init() tea.Cmd { return nil }

func (m CrawlerFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	if !isKey {
		return m, nil
	}
	k := key.String()

	if k == "ctrl+c" {
		m.quit = true
		return m, tea.Quit
	}

	// Confirm stage
	if m.stage == cStageConfirm {
		switch k {
		case "enter", "y":
			m.stage = cStageDone
			return m, tea.Quit
		case "n", "esc", "backspace":
			m.stage = cStageForm
			m.focusField(cfSubmit)
		}
		return m, nil
	}

	// Form stage
	if k == "esc" && !m.isTextFocused() {
		m.quit = true
		return m, tea.Quit
	}

	if k == " " {
		if m.isTextFocused() {
			m.routeRune(' ')
		}
		return m, nil
	}

	switch k {
	case "tab", "down":
		m.focusField(m.nextField())

	case "shift+tab", "up":
		m.focusField(m.prevField())

	case "enter":
		switch m.focus {
		case cfSubmit:
			m.errMsg = m.validate()
			if m.errMsg == "" {
				m.stage = cStageConfirm
			}
		case cfURLs:
			m.urlsInput.Commit()
		default:
			m.focusField(m.nextField())
		}

	case "left":
		switch m.focus {
		case cfDepth:
			m.depthStep.Dec()
		case cfMaxPages:
			m.maxPagesStep.Dec()
		case cfDelay:
			m.delayStep.Dec()
		default:
			m.routeLeft()
		}

	case "right":
		switch m.focus {
		case cfDepth:
			m.depthStep.Inc()
		case cfMaxPages:
			m.maxPagesStep.Inc()
		case cfDelay:
			m.delayStep.Inc()
		default:
			m.routeRight()
		}

	case "backspace":
		m.routeBackspace()

	default:
		if len(key.Runes) == 1 {
			m.routeRune(key.Runes[0])
		}
	}

	return m, nil
}

func (m *CrawlerFormModel) routeRune(r rune) {
	switch m.focus {
	case cfURLs:
		m.urlsInput.HandleRune(r)
	case cfFile:
		m.fileInput.HandleRune(r)
	case cfJSON:
		m.jsonInput.HandleRune(r)
	case cfStorage:
		m.storageInput.HandleRune(r)
	}
}

func (m *CrawlerFormModel) routeBackspace() {
	switch m.focus {
	case cfURLs:
		m.urlsInput.HandleBackspace()
	case cfFile:
		m.fileInput.HandleBackspace()
	case cfJSON:
		m.jsonInput.HandleBackspace()
	case cfStorage:
		m.storageInput.HandleBackspace()
	}
}

func (m *CrawlerFormModel) routeLeft() {
	switch m.focus {
	case cfURLs:
		m.urlsInput.HandleLeft()
	case cfFile:
		m.fileInput.HandleLeft()
	case cfJSON:
		m.jsonInput.HandleLeft()
	case cfStorage:
		m.storageInput.HandleLeft()
	}
}

func (m *CrawlerFormModel) routeRight() {
	switch m.focus {
	case cfURLs:
		m.urlsInput.HandleRight()
	case cfFile:
		m.fileInput.HandleRight()
	case cfJSON:
		m.jsonInput.HandleRight()
	case cfStorage:
		m.storageInput.HandleRight()
	}
}

func (m CrawlerFormModel) validate() string {
	hasURL := len(m.urlsInput.Values) > 0 || strings.TrimSpace(m.urlsInput.current.Value) != ""
	hasFile := strings.TrimSpace(m.fileInput.Value) != ""
	hasJSON := strings.TrimSpace(m.jsonInput.Value) != ""
	if !hasURL && !hasFile && !hasJSON {
		return "Provide at least one URL, file, or JSON source"
	}
	return ""
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m CrawlerFormModel) View() string {
	var page strings.Builder

	page.WriteString(renderCrawlerForm(m))
	page.WriteString("\n")

	if m.stage == cStageConfirm {
		page.WriteString(renderCrawlerConfirmBox(m.Result()))
		page.WriteString("\n\n")
		page.WriteString(
			renderKeyHints("enter / y", "yes, start crawl", "n / esc", "go back and edit", "ctrl+c", "quit"),
		)
		page.WriteString("\n")
	}

	return page.String()
}

func renderCrawlerForm(m CrawlerFormModel) string {
	var sb strings.Builder
	locked := m.stage == cStageConfirm

	if locked {
		sb.WriteString(styleFormTitleLocked.Render("  Crawler Configuration  ✔") + "\n")
	} else {
		sb.WriteString(styleFormTitle.Render("  Crawler Configuration") + "\n")
	}
	sb.WriteString(divider() + "\n\n")

	sb.WriteString(m.urlsInput.View() + "\n\n")
	sb.WriteString(m.fileInput.View() + "\n\n")
	sb.WriteString(m.jsonInput.View() + "\n\n")
	sb.WriteString(m.depthStep.View() + "\n\n")
	sb.WriteString(m.maxPagesStep.View() + "\n\n")
	sb.WriteString(m.delayStep.View() + "\n\n")
	sb.WriteString(m.storageInput.View() + "\n\n")

	sb.WriteString(divider() + "\n  ")
	switch {
	case locked:
		sb.WriteString(styleSubmitDone.Render("  ✔   SUBMITTED  "))
	case m.focus == cfSubmit:
		sb.WriteString(styleSubmitActive.Render("  ▶   START CRAWL  "))
	default:
		sb.WriteString(styleSubmitInactive.Render("  ▶   START CRAWL  "))
	}
	sb.WriteString("\n")

	if m.errMsg != "" {
		sb.WriteString("\n" + styleError.Render("  ✘  "+m.errMsg) + "\n")
	}

	if !locked {
		sb.WriteString("\n")
		sb.WriteString(renderKeyHints("tab", "next", "shift+tab", "prev", "← →", "change") + "\n")
		sb.WriteString(renderKeyHints("enter", "confirm/add URL", "ctrl+c", "quit"))
	}

	box := styleFormBox
	if locked {
		box = styleFormBoxLocked
	}
	return box.Render(sb.String())
}

func renderCrawlerConfirmBox(cfg CrawlerConfig) string {
	var sb strings.Builder

	keyStyle := lipgloss.NewStyle().Foreground(clrGray).Width(12).Align(lipgloss.Right)
	valStyle := lipgloss.NewStyle().Foreground(clrWhite).Bold(true)
	sep := lipgloss.NewStyle().Foreground(clrGray).Render("  │  ")

	row := func(k, v string) {
		sb.WriteString(keyStyle.Render(k) + sep + valStyle.Render(v) + "\n")
	}

	sb.WriteString(styleConfirmTitle.Render("⚙  Confirm — ready to crawl?") + "\n")
	sb.WriteString(styleDim.Render(strings.Repeat("─", 46)) + "\n")

	for i, u := range cfg.URLs {
		label := "URLs"
		if i > 0 {
			label = ""
		}
		row(label, u)
	}
	if len(cfg.URLs) == 0 {
		row("URLs", "(from file/json)")
	}
	row("File", strOrDash(cfg.File))
	row("JSON", strOrDash(cfg.JSON))
	row("Depth", fmt.Sprintf("%d", cfg.Depth))
	row("Max Pages", fmt.Sprintf("%d", cfg.MaxPages))
	row("Delay (s)", fmt.Sprintf("%d", cfg.Delay))
	row("Storage", cfg.Storage)

	return styleConfirmBox.Render(sb.String())
}

// ── Result + Runner ───────────────────────────────────────────────────────────

func (m CrawlerFormModel) Result() CrawlerConfig {
	urls := append([]string{}, m.urlsInput.Values...)
	if v := strings.TrimSpace(m.urlsInput.current.Value); v != "" {
		urls = append(urls, v)
	}
	storage := m.storageInput.Value
	if storage == "" {
		storage = "data/pages"
	}
	return CrawlerConfig{
		URLs:      urls,
		File:      strings.TrimSpace(m.fileInput.Value),
		JSON:      strings.TrimSpace(m.jsonInput.Value),
		Depth:     m.depthStep.Value,
		MaxPages:  m.maxPagesStep.Value,
		Delay:     m.delayStep.Value,
		Storage:   storage,
		Submitted: m.stage == cStageDone,
	}
}

// RunCrawlerForm runs inline — banner stays visible above, results print below.
func RunCrawlerForm(defaults CrawlerConfig) (CrawlerConfig, error) {
	p := tea.NewProgram(NewCrawlerForm(defaults)) // no WithAltScreen
	final, err := p.Run()
	if err != nil {
		return CrawlerConfig{}, err
	}
	return final.(CrawlerFormModel).Result(), nil
}

func strOrDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
