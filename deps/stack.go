package deps

type stack struct {
	data []string
}

func newVendorStack(values ...string) *stack {
	limit := len(values) * 2
	s := &stack{data: make([]string, 0, limit)}
	for i := len(values) - 1; i >= 0; i-- {
		s.push(values[i])
	}
	return s
}

func (s stack) empty() bool {
	return s.len() == 0
}

func (s stack) len() int {
	return len(s.data)
}

func (s *stack) push(values ...string) {
	s.data = append(s.data, values...)
}

func (s *stack) pop() string {
	if s.len() > 0 {
		last := s.data[s.len()-1]
		s.data = s.data[:s.len()-1]
		return last
	}
	return ""
}
