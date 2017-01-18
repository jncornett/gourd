package main

type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

type BufferedScanner interface {
	Scanner
	Unscan()
}

type bufferedScanner struct {
	s        Scanner
	last     string
	buffered bool
}

func NewBufferedScanner(s Scanner) BufferedScanner {
	return &bufferedScanner{s: s}
}

func (s *bufferedScanner) Scan() bool {
	if s.buffered {
		return true
	} /* else { s.buffered = false } */
	return s.s.Scan()
}

func (s *bufferedScanner) Text() string {
	if s.buffered {
		s.buffered = false
	} else {
		s.last = s.s.Text()
	}
	return s.last
}

func (s *bufferedScanner) Err() error {
	return s.s.Err()
}

func (s *bufferedScanner) Unscan() {
	s.buffered = true
}

type filteredScanner struct {
	s      Scanner
	filter func(string) bool
}

func NewFilteredScanner(s Scanner, filter func(string) bool) Scanner {
	return &filteredScanner{s: s, filter: filter}
}

func (s *filteredScanner) Scan() bool {
	for s.s.Scan() {
		// FIXME calling filteredScanner.s.Text allocates a new string buffer
		if s.filter(s.s.Text()) {
			return true
		}
	}
	return false
}

func (s *filteredScanner) Text() string {
	return s.s.Text()
}

func (s *filteredScanner) Err() error {
	return s.s.Err()
}
