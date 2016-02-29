// Copyright 2016 github.com/govend/govend. All rights reserved.
// Use of this source code is governed by an Apache 2.0
// license that can be found in the LICENSE file.

package deps

type dep struct {
	path  string
	level int
}

type stack struct {
	data []dep
}

func newVendorStack(values ...string) *stack {
	limit := len(values) * 2
	s := &stack{data: make([]dep, 0, limit)}
	for i := len(values) - 1; i >= 0; i-- {
		s.push(0, values[i])
	}
	return s
}

func (s stack) empty() bool {
	return s.len() == 0
}

func (s stack) len() int {
	return len(s.data)
}

func (s *stack) push(level int, values ...string) {
	for _, val := range values {
		s.data = append(s.data, dep{path: val, level: level})
	}
}

func (s *stack) pop() dep {
	if s.len() > 0 {
		last := s.data[s.len()-1]
		s.data = s.data[:s.len()-1]
		return last
	}
	return dep{}
}
