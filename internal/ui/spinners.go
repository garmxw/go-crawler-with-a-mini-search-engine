package ui

import "github.com/pterm/pterm"

// Spinner wraps pterm's SpinnerPrinter for clean start/stop API.
// Sequence uses only plain ASCII chars (- \ | /) so it renders correctly
// on Windows Terminal Preview and the legacy conhost fallback.
type Spinner struct {
	sp *pterm.SpinnerPrinter
}

// NewSpinner starts a spinner with the given message. Call Done or Fail to stop it.
func NewSpinner(msg string) *Spinner {
	sp, _ := pterm.DefaultSpinner.
		WithSequence("-", "\\", "|", "/"). // ASCII-only: safe on all Windows consoles
		WithStyle(&pterm.Style{pterm.FgCyan}).
		WithMessageStyle(&pterm.Style{pterm.FgWhite}).
		Start(msg)

	return &Spinner{sp: sp}
}

// UpdateMsg changes the spinner text while it's running.
func (s *Spinner) UpdateMsg(msg string) {
	if s.sp != nil {
		s.sp.UpdateText(msg)
	}
}

// Done stops the spinner with a green success message.
func (s *Spinner) Done(msg string) {
	if s.sp != nil {
		s.sp.Success(msg)
	}
}

// Fail stops the spinner with a red failure message.
func (s *Spinner) Fail(msg string) {
	if s.sp != nil {
		s.sp.Fail(msg)
	}
}

// Stop silently stops the spinner.
func (s *Spinner) Stop() {
	if s.sp != nil {
		s.sp.Stop()
	}
}
