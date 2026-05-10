package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// tea is used for tea.Model, tea.Cmd, tea.KeyMsg, tea.Quit in Update().
// newProgram() (defined in theme.go) wraps tea.NewProgram with Windows-safe options.

// ── Config ────────────────────────────────────────────────────────────────────

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
	Delay     int // crawl delay in seconds (live mode), default 2
	Storage   string
	Submitted bool
}

// ── Stage ─────────────────────────────────────────────────────────────────────

type searchStage int

const (
	stageForm    searchStage = iota // user filling the form
	stageConfirm                    // confirm panel visible below form, waiting for y/enter
	stageDone                       // user confirmed — program exits
)

// ── Field indices ─────────────────────────────────────────────────────────────

const (
	sfMode = iota
	sfQuery
	sfLang
	sfLimit
	sfDetailed
	sfPath
	sfURLs
	sfFile
	sfJSON
	sfDepth
	sfMaxPages
	sfDelay
	sfStorage
	sfSubmit
)

// ── Model ─────────────────────────────────────────────────────────────────────

type SearchFormModel struct {
	modeSelect   Select
	queryInput   TextInput
	langSelect   Select
	limitStep    Stepper
	detailedTog  Toggle
	pathInput    TextInput
	urlsInput    MultiInput
	fileInput    TextInput
	jsonInput    TextInput
	depthStep    Stepper
	maxPagesStep Stepper
	delayStep    Stepper
	storageInput TextInput

	focus  int
	stage  searchStage
	errMsg string
	quit   bool
}

func NewSearchForm(d SearchConfig) SearchFormModel {
	modeIdx := 0
	for i, o := range []string{"local", "web", "live"} {
		if o == d.Mode {
			modeIdx = i
		}
	}
	langIdx := 0
	if d.Lang == "french" {
		langIdx = 1
	}
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

	m := SearchFormModel{
		modeSelect:   Select{Label: "Mode", Options: []string{"local", "web", "live"}, Index: modeIdx},
		queryInput:   NewTextInput("Query", "type your search query…"),
		langSelect:   Select{Label: "Language", Options: []string{"english", "french"}, Index: langIdx},
		limitStep:    NewStepper("Limit", d.Limit, 1, 100),
		detailedTog:  NewToggle("Detailed", d.Detailed),
		pathInput:    NewTextInput("Local Path", "path/to/docs"),
		urlsInput:    NewMultiInput("URLs", "https://example.com"),
		fileInput:    NewTextInput("URL File", "urls.txt"),
		jsonInput:    NewTextInput("JSON File", "urls.json"),
		depthStep:    NewStepper("Depth", d.Depth, 0, 10),
		maxPagesStep: NewStepper("Max Pages", d.MaxPages, 1, -1),
		delayStep:    NewStepper("Delay (s)", d.Delay, 0, 60),
		storageInput: NewTextInput("Storage", "data/pages"),
	}

	m.queryInput.Value = d.Query
	m.queryInput.cursor = len([]rune(d.Query))
	m.pathInput.Value = d.Path
	m.pathInput.cursor = len([]rune(d.Path))
	m.fileInput.Value = d.File
	m.fileInput.cursor = len([]rune(d.File))
	m.jsonInput.Value = d.JSON
	m.jsonInput.cursor = len([]rune(d.JSON))
	m.storageInput.Value = d.Storage
	m.storageInput.cursor = len([]rune(d.Storage))
	m.urlsInput.Values = d.URLs

	m.focusField(sfMode)
	return m
}

func (m SearchFormModel) isTextFocused() bool {
	switch m.focus {
	case sfQuery, sfPath, sfURLs, sfFile, sfJSON, sfStorage:
		return true
	}
	return false
}

func (m *SearchFormModel) focusField(idx int) {
	m.modeSelect.Blur()
	m.queryInput.Blur()
	m.langSelect.Blur()
	m.limitStep.Blur()
	m.detailedTog.Blur()
	m.pathInput.Blur()
	m.urlsInput.Blur()
	m.fileInput.Blur()
	m.jsonInput.Blur()
	m.depthStep.Blur()
	m.maxPagesStep.Blur()
	m.delayStep.Blur()
	m.storageInput.Blur()

	m.focus = idx
	switch idx {
	case sfMode:
		m.modeSelect.Focus()
	case sfQuery:
		m.queryInput.Focus()
	case sfLang:
		m.langSelect.Focus()
	case sfLimit:
		m.limitStep.Focus()
	case sfDetailed:
		m.detailedTog.Focus()
	case sfPath:
		m.pathInput.Focus()
	case sfURLs:
		m.urlsInput.Focus()
	case sfFile:
		m.fileInput.Focus()
	case sfJSON:
		m.jsonInput.Focus()
	case sfDepth:
		m.depthStep.Focus()
	case sfMaxPages:
		m.maxPagesStep.Focus()
	case sfDelay:
		m.delayStep.Focus()
	case sfStorage:
		m.storageInput.Focus()
	}
}

// visibleFields returns field indices to render, conditioned on mode.
func (m SearchFormModel) visibleFields() []int {
	base := []int{sfMode, sfQuery, sfLang, sfLimit, sfDetailed}
	switch m.modeSelect.Value() {
	case "local":
		base = append(base, sfPath)
	case "web":
		base = append(base, sfStorage)
	case "live":
		base = append(base, sfURLs, sfFile, sfJSON, sfDepth, sfMaxPages, sfDelay, sfStorage)
	}
	return append(base, sfSubmit)
}

func (m SearchFormModel) nextField() int {
	fields := m.visibleFields()
	for i, f := range fields {
		if f == m.focus {
			if i+1 < len(fields) {
				return fields[i+1]
			}
			return fields[0]
		}
	}
	return fields[0]
}

func (m SearchFormModel) prevField() int {
	fields := m.visibleFields()
	for i, f := range fields {
		if f == m.focus {
			if i > 0 {
				return fields[i-1]
			}
			return fields[len(fields)-1]
		}
	}
	return fields[len(fields)-1]
}

// ── Bubbletea ─────────────────────────────────────────────────────────────────

func (m SearchFormModel) Init() tea.Cmd { return nil }

func (m SearchFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, isKey := msg.(tea.KeyMsg)
	if !isKey {
		return m, nil
	}
	k := key.String()

	// Always handle quit
	if k == "ctrl+c" {
		m.quit = true
		return m, tea.Quit
	}

	// ── Confirm stage: only waiting for y/enter or n/esc to go back ───────────
	if m.stage == stageConfirm {
		switch k {
		case "enter", "y":
			m.stage = stageDone
			return m, tea.Quit
		case "n", "esc", "backspace":
			m.stage = stageForm
			m.focusField(sfSubmit)
		}
		return m, nil
	}

	// ── Form stage ────────────────────────────────────────────────────────────

	// esc quits only when not in a text field
	if k == "esc" && !m.isTextFocused() {
		m.quit = true
		return m, tea.Quit
	}

	// Space: explicit routing — text fields get a literal space, toggle flips
	if k == " " {
		if m.focus == sfDetailed {
			m.detailedTog.Flip()
		} else if m.isTextFocused() {
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
		case sfSubmit:
			m.errMsg = m.validate()
			if m.errMsg == "" {
				// Lock the form and show confirm panel below it
				m.stage = stageConfirm
			}
		case sfURLs:
			m.urlsInput.Commit()
		default:
			m.focusField(m.nextField())
		}

	case "left":
		switch m.focus {
		case sfMode:
			m.modeSelect.Prev()
		case sfLang:
			m.langSelect.Prev()
		case sfLimit:
			m.limitStep.Dec()
		case sfDepth:
			m.depthStep.Dec()
		case sfMaxPages:
			m.maxPagesStep.Dec()
		case sfDelay:
			m.delayStep.Dec()
		default:
			m.routeLeft()
		}

	case "right":
		switch m.focus {
		case sfMode:
			m.modeSelect.Next()
		case sfLang:
			m.langSelect.Next()
		case sfLimit:
			m.limitStep.Inc()
		case sfDepth:
			m.depthStep.Inc()
		case sfMaxPages:
			m.maxPagesStep.Inc()
		case sfDelay:
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

func (m *SearchFormModel) routeRune(r rune) {
	switch m.focus {
	case sfQuery:
		m.queryInput.HandleRune(r)
	case sfPath:
		m.pathInput.HandleRune(r)
	case sfURLs:
		m.urlsInput.HandleRune(r)
	case sfFile:
		m.fileInput.HandleRune(r)
	case sfJSON:
		m.jsonInput.HandleRune(r)
	case sfStorage:
		m.storageInput.HandleRune(r)
	}
}

func (m *SearchFormModel) routeBackspace() {
	switch m.focus {
	case sfQuery:
		m.queryInput.HandleBackspace()
	case sfPath:
		m.pathInput.HandleBackspace()
	case sfURLs:
		m.urlsInput.HandleBackspace()
	case sfFile:
		m.fileInput.HandleBackspace()
	case sfJSON:
		m.jsonInput.HandleBackspace()
	case sfStorage:
		m.storageInput.HandleBackspace()
	}
}

func (m *SearchFormModel) routeLeft() {
	switch m.focus {
	case sfQuery:
		m.queryInput.HandleLeft()
	case sfPath:
		m.pathInput.HandleLeft()
	case sfURLs:
		m.urlsInput.HandleLeft()
	case sfFile:
		m.fileInput.HandleLeft()
	case sfJSON:
		m.jsonInput.HandleLeft()
	case sfStorage:
		m.storageInput.HandleLeft()
	}
}

func (m *SearchFormModel) routeRight() {
	switch m.focus {
	case sfQuery:
		m.queryInput.HandleRight()
	case sfPath:
		m.pathInput.HandleRight()
	case sfURLs:
		m.urlsInput.HandleRight()
	case sfFile:
		m.fileInput.HandleRight()
	case sfJSON:
		m.jsonInput.HandleRight()
	case sfStorage:
		m.storageInput.HandleRight()
	}
}

func (m SearchFormModel) validate() string {
	if strings.TrimSpace(m.queryInput.Value) == "" {
		return "Query cannot be empty"
	}
	if m.modeSelect.Value() == "local" && strings.TrimSpace(m.pathInput.Value) == "" {
		return "Local mode requires a path"
	}
	if m.modeSelect.Value() == "live" &&
		len(m.urlsInput.Values) == 0 &&
		strings.TrimSpace(m.urlsInput.current.Value) == "" &&
		strings.TrimSpace(m.fileInput.Value) == "" &&
		strings.TrimSpace(m.jsonInput.Value) == "" {
		return "Live mode requires at least one URL, file, or JSON source"
	}
	return ""
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m SearchFormModel) View() string {
	var page strings.Builder

	// Form is always visible
	page.WriteString(renderSearchForm(m))
	page.WriteString("\n")

	// Confirm panel appears below the locked form
	if m.stage == stageConfirm {
		page.WriteString(renderSearchConfirmBox(m.Result()))
		page.WriteString("\n\n")
		page.WriteString(
			renderKeyHints("enter / y", "yes, run it", "n / esc", "go back and edit", "ctrl+c", "quit"),
		)
		page.WriteString("\n")
	}

	return page.String()
}

// renderSearchForm renders the form box itself.
func renderSearchForm(m SearchFormModel) string {
	var sb strings.Builder
	locked := m.stage == stageConfirm

	if locked {
		sb.WriteString(styleFormTitleLocked.Render("  Search Configuration  ✔") + "\n")
	} else {
		sb.WriteString(styleFormTitle.Render("  Search Configuration") + "\n")
	}
	sb.WriteString(divider() + "\n\n")

	for _, f := range m.visibleFields() {
		switch f {
		case sfMode:
			sb.WriteString(m.modeSelect.View() + "\n\n")
		case sfQuery:
			sb.WriteString(m.queryInput.View() + "\n\n")
		case sfLang:
			sb.WriteString(m.langSelect.View() + "\n\n")
		case sfLimit:
			sb.WriteString(m.limitStep.View() + "\n\n")
		case sfDetailed:
			sb.WriteString(m.detailedTog.View() + "\n\n")
		case sfPath:
			sb.WriteString(m.pathInput.View() + "\n\n")
		case sfURLs:
			sb.WriteString(m.urlsInput.View() + "\n\n")
		case sfFile:
			sb.WriteString(m.fileInput.View() + "\n\n")
		case sfJSON:
			sb.WriteString(m.jsonInput.View() + "\n\n")
		case sfDepth:
			sb.WriteString(m.depthStep.View() + "\n\n")
		case sfMaxPages:
			sb.WriteString(m.maxPagesStep.View() + "\n\n")
		case sfDelay:
			sb.WriteString(m.delayStep.View() + "\n\n")
		case sfStorage:
			sb.WriteString(m.storageInput.View() + "\n\n")
		case sfSubmit:
			sb.WriteString(divider() + "\n  ")
			switch {
			case locked:
				sb.WriteString(styleSubmitDone.Render("  ✔   SUBMITTED  "))
			case m.focus == sfSubmit:
				sb.WriteString(styleSubmitActive.Render("  ▶   RUN SEARCH  "))
			default:
				sb.WriteString(styleSubmitInactive.Render("  ▶   RUN SEARCH  "))
			}
			sb.WriteString("\n")
		}
	}

	if m.errMsg != "" {
		sb.WriteString("\n" + styleError.Render("  ✘  "+m.errMsg) + "\n")
	}

	if !locked {
		sb.WriteString("\n")
		sb.WriteString(renderKeyHints("tab", "next", "shift+tab", "prev", "← →", "change") + "\n")
		sb.WriteString(renderKeyHints("enter", "confirm/add", "space", "toggle/type", "ctrl+c", "quit"))
	}

	box := styleFormBox
	if locked {
		box = styleFormBoxLocked
	}
	return box.Render(sb.String())
}

// renderSearchConfirmBox renders the summary table that appears below the form.
func renderSearchConfirmBox(cfg SearchConfig) string {
	var sb strings.Builder

	keyStyle := lipgloss.NewStyle().Foreground(clrGray).Width(12).Align(lipgloss.Right)
	valStyle := lipgloss.NewStyle().Foreground(clrWhite).Bold(true)
	sep := lipgloss.NewStyle().Foreground(clrGray).Render("  │  ")

	row := func(k, v string) {
		sb.WriteString(keyStyle.Render(k) + sep + valStyle.Render(v) + "\n")
	}

	modeBg := clrCyan
	switch cfg.Mode {
	case "local":
		modeBg = clrGreen
	case "live":
		modeBg = clrOrange
	}
	badge := lipgloss.NewStyle().Background(modeBg).Foreground(clrDark).Bold(true).Padding(0, 1)

	sb.WriteString(styleConfirmTitle.Render("⚙  Confirm — ready to run?") + "\n")
	sb.WriteString(styleDim.Render(strings.Repeat("─", 46)) + "\n")
	sb.WriteString(keyStyle.Render("Mode") + sep + badge.Render(" "+cfg.Mode+" ") + "\n")
	row("Query", fmt.Sprintf("%q", cfg.Query))
	row("Language", cfg.Lang)
	row("Limit", fmt.Sprintf("%d", cfg.Limit))
	row("Detailed", fmt.Sprintf("%v", cfg.Detailed))

	switch cfg.Mode {
	case "local":
		row("Path", cfg.Path)
	case "web":
		row("Storage", cfg.Storage)
	case "live":
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
		row("Depth", fmt.Sprintf("%d", cfg.Depth))
		row("Max Pages", fmt.Sprintf("%d", cfg.MaxPages))
		row("Delay (s)", fmt.Sprintf("%d", cfg.Delay))
		row("Storage", cfg.Storage)
	}

	return styleConfirmBox.Render(sb.String())
}

// ── Result + Runner ───────────────────────────────────────────────────────────

func (m SearchFormModel) Result() SearchConfig {
	urls := append([]string{}, m.urlsInput.Values...)
	if v := strings.TrimSpace(m.urlsInput.current.Value); v != "" {
		urls = append(urls, v)
	}
	storage := m.storageInput.Value
	if storage == "" {
		storage = "data/pages"
	}
	return SearchConfig{
		Mode:      m.modeSelect.Value(),
		Query:     strings.TrimSpace(m.queryInput.Value),
		Path:      strings.TrimSpace(m.pathInput.Value),
		URLs:      urls,
		Lang:      m.langSelect.Value(),
		Limit:     m.limitStep.Value,
		Detailed:  m.detailedTog.Value,
		File:      strings.TrimSpace(m.fileInput.Value),
		JSON:      strings.TrimSpace(m.jsonInput.Value),
		Depth:     m.depthStep.Value,
		MaxPages:  m.maxPagesStep.Value,
		Delay:     m.delayStep.Value,
		Storage:   storage,
		Submitted: m.stage == stageDone,
	}
}

// RunSearchForm runs the form inline (no alt screen).
// Uses newProgram() which sets tea.WithInput(os.Stdin) + tea.WithOutput(os.Stderr)
// so it works correctly on Windows Terminal Preview and all Unix terminals.
func RunSearchForm(defaults SearchConfig) (SearchConfig, error) {
	p := newProgram(NewSearchForm(defaults))
	final, err := p.Run()
	restoreConsoleInput() // restore stdin mode on Windows before crawler runs
	if err != nil {
		return SearchConfig{}, err
	}
	return final.(SearchFormModel).Result(), nil
}
