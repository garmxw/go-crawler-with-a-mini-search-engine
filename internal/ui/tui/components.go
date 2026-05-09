package tui

import (
	"fmt"
	"strings"
)

const labelWidth = 16

func labelPad() string { return strings.Repeat(" ", labelWidth) }

// ── TextInput ─────────────────────────────────────────────────────────────────

type TextInput struct {
	Label       string
	Value       string
	Placeholder string
	cursor      int
	focused     bool
}

func NewTextInput(label, placeholder string) TextInput {
	return TextInput{Label: label, Placeholder: placeholder}
}

func (t *TextInput) Focus()        { t.focused = true }
func (t *TextInput) Blur()         { t.focused = false }
func (t *TextInput) Focused() bool { return t.focused }

func (t *TextInput) HandleRune(r rune) {
	runes := []rune(t.Value)
	out := make([]rune, len(runes)+1)
	copy(out, runes[:t.cursor])
	out[t.cursor] = r
	copy(out[t.cursor+1:], runes[t.cursor:])
	t.Value = string(out)
	t.cursor++
}

func (t *TextInput) HandleBackspace() {
	if t.cursor == 0 {
		return
	}
	runes := []rune(t.Value)
	runes = append(runes[:t.cursor-1], runes[t.cursor:]...)
	t.Value = string(runes)
	t.cursor--
}

func (t *TextInput) HandleLeft() {
	if t.cursor > 0 {
		t.cursor--
	}
}
func (t *TextInput) HandleRight() {
	if t.cursor < len([]rune(t.Value)) {
		t.cursor++
	}
}

func (t TextInput) View() string {
	label := styleUnfocusedLabel.Render(t.Label)
	if t.focused {
		label = styleFocusedLabel.Render(t.Label)
	}
	var content string
	if t.focused {
		runes := []rune(t.Value)
		before := string(runes[:t.cursor])
		cursorChar := " "
		after := ""
		if t.cursor < len(runes) {
			cursorChar = string(runes[t.cursor])
			after = string(runes[t.cursor+1:])
		}
		inner := before + styleCursor.Render(cursorChar) + after
		if t.Value == "" && t.Placeholder != "" {
			inner = styleCursor.Render(" ") + styleHint.Render(t.Placeholder)
		}
		content = styleFocusedInput.Render(inner)
	} else {
		display := t.Value
		if display == "" && t.Placeholder != "" {
			display = styleHint.Render(t.Placeholder)
		}
		content = styleUnfocusedInput.Render(display)
	}
	return label + content
}

// ── Select ────────────────────────────────────────────────────────────────────

type Select struct {
	Label   string
	Options []string
	Index   int
	focused bool
}

func NewSelect(label string, options []string) Select {
	return Select{Label: label, Options: options}
}

func (s *Select) Focus()        { s.focused = true }
func (s *Select) Blur()         { s.focused = false }
func (s *Select) Focused() bool { return s.focused }
func (s *Select) Value() string { return s.Options[s.Index] }
func (s *Select) Next()         { s.Index = (s.Index + 1) % len(s.Options) }
func (s *Select) Prev()         { s.Index = (s.Index - 1 + len(s.Options)) % len(s.Options) }

func (s Select) View() string {
	label := styleUnfocusedLabel.Render(s.Label)
	if s.focused {
		label = styleFocusedLabel.Render(s.Label)
	}
	var badges []string
	for i, opt := range s.Options {
		if i == s.Index {
			badges = append(badges, styleSelectedBadge.Render(" "+opt+" "))
		} else {
			badges = append(badges, styleUnselectedBadge.Render(" "+opt+" "))
		}
	}
	hint := ""
	if s.focused {
		hint = styleHint.Render("   ← →")
	}
	return label + strings.Join(badges, " ") + hint
}

// ── Stepper ───────────────────────────────────────────────────────────────────

type Stepper struct {
	Label   string
	Value   int
	Min     int
	Max     int // -1 = no limit
	focused bool
}

func NewStepper(label string, value, min, max int) Stepper {
	return Stepper{Label: label, Value: value, Min: min, Max: max}
}

func (s *Stepper) Focus()        { s.focused = true }
func (s *Stepper) Blur()         { s.focused = false }
func (s *Stepper) Focused() bool { return s.focused }
func (s *Stepper) Inc() {
	if s.Max < 0 || s.Value < s.Max {
		s.Value++
	}
}
func (s *Stepper) Dec() {
	if s.Value > s.Min {
		s.Value--
	}
}

func (s Stepper) View() string {
	label := styleUnfocusedLabel.Render(s.Label)
	if s.focused {
		label = styleFocusedLabel.Render(s.Label)
	}
	dec := styleUnselectedBadge.Render("  ◀  ")
	inc := styleUnselectedBadge.Render("  ▶  ")
	if s.focused {
		dec = styleSelectedBadge.Render("  ◀  ")
		inc = styleSelectedBadge.Render("  ▶  ")
	}
	val := styleFocusedLabel.Render(fmt.Sprintf("  %d  ", s.Value))
	hint := ""
	if s.focused {
		hint = styleHint.Render("   ← →")
	}
	return label + dec + val + inc + hint
}

// ── Toggle ────────────────────────────────────────────────────────────────────

type Toggle struct {
	Label   string
	Value   bool
	focused bool
}

func NewToggle(label string, value bool) Toggle {
	return Toggle{Label: label, Value: value}
}

func (t *Toggle) Focus()        { t.focused = true }
func (t *Toggle) Blur()         { t.focused = false }
func (t *Toggle) Focused() bool { return t.focused }
func (t *Toggle) Flip()         { t.Value = !t.Value }

func (t Toggle) View() string {
	label := styleUnfocusedLabel.Render(t.Label)
	if t.focused {
		label = styleFocusedLabel.Render(t.Label)
	}
	offStyle := styleUnselectedBadge
	onStyle := styleUnselectedBadge
	if t.Value {
		onStyle = styleOnBadge
	} else {
		offStyle = styleOffBadge
	}
	hint := ""
	if t.focused {
		hint = styleHint.Render("   space to toggle")
	}
	sep := styleDim.Render("  │  ")
	return label + offStyle.Render(" OFF ") + sep + onStyle.Render(" ON ") + hint
}

// ── MultiInput ────────────────────────────────────────────────────────────────

type MultiInput struct {
	Label   string
	Values  []string
	current TextInput
	focused bool
}

func NewMultiInput(label, placeholder string) MultiInput {
	return MultiInput{
		Label:   label,
		current: NewTextInput("", placeholder),
	}
}

func (m *MultiInput) Focus() {
	m.focused = true
	m.current.Focus()
}
func (m *MultiInput) Blur() {
	m.focused = false
	m.current.Blur()
}
func (m *MultiInput) Focused() bool { return m.focused }

func (m *MultiInput) Commit() {
	v := strings.TrimSpace(m.current.Value)
	if v != "" {
		m.Values = append(m.Values, v)
		m.current.Value = ""
		m.current.cursor = 0
	}
}

func (m *MultiInput) DeleteLast() {
	if len(m.Values) > 0 {
		m.Values = m.Values[:len(m.Values)-1]
	}
}

func (m *MultiInput) HandleRune(r rune)  { m.current.HandleRune(r) }
func (m *MultiInput) HandleLeft()        { m.current.HandleLeft() }
func (m *MultiInput) HandleRight()       { m.current.HandleRight() }

func (m *MultiInput) HandleBackspace() {
	if m.current.Value == "" {
		m.DeleteLast()
	} else {
		m.current.HandleBackspace()
	}
}

func (m MultiInput) View() string {
	label := styleUnfocusedLabel.Render(m.Label)
	if m.focused {
		label = styleFocusedLabel.Render(m.Label)
	}

	var lines []string
	lines = append(lines, label+m.current.View())

	pad := labelPad()
	for i, v := range m.Values {
		num := styleDim.Render(fmt.Sprintf("  %d.  ", i+1))
		tag := styleURLTag.Render(" "+v+" ")
		lines = append(lines, pad+num+tag)
	}

	if m.focused {
		hint := styleHint.Render("  enter ↵ add url   backspace on empty removes last")
		lines = append(lines, pad+hint)
	}

	return strings.Join(lines, "\n")
}
