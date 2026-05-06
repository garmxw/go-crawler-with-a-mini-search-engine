package cli

// custom flag type for multi-value flags
type MultiFlag []string

func (m *MultiFlag) String() string {
	return ""
}

func (m *MultiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}
