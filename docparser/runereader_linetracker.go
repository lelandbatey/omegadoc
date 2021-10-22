package docparser

import (
	"bufio"
)

// runereader_linetracker.go holds a struct which wraps a bufio.Reader,
// exposing only the "ReadRune" method and a new method, the "LineNumber"
// method, which tells you the line-number of the last-read rune in the
// underlying bufio.Reader. We use this so we can keep track of the
// line-numbers where things are defined, whether those are errors or important
// identifiers like the start of an OmegadDoc.

type linetracker struct {
	brdr   *bufio.Reader
	lineno int
}

func (lt linetracker) ReadRune() (rune, int, error) {
	r, s, err := lt.brdr.ReadRune()
	if r == '\n' {
		lt.lineno++
	}
	return r, s, err
}

func (lt linetracker) LineNumber() int {
	return lt.lineno
}
