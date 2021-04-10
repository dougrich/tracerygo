package tracerygo

import (
	"errors"
	"io"
	"strings"
)

type ModifierFunc func(io.Writer) Modifier

type Modifier interface {
	io.Writer
	Finalize()
}

var (
	ErrUnexpectedNumberOfBytesWritten = errors.New("Unexpected number of bytes written")
)

type capitalizePipe struct {
	out  io.Writer
	done bool
}

func ModifierCapitalize(out io.Writer) Modifier {
	return &capitalizePipe{out: out, done: false}
}

func (p *capitalizePipe) Write(b []byte) (int, error) {
	if p.done {
		return p.out.Write(b)
	} else {
		p.done = true
		s := string(b)
		s = strings.ToUpper(string(s[0])) + s[1:]
		l, err := p.out.Write([]byte(s))
		if l == len(s) {
			return len(b), err
		} else {
			return 0, ErrUnexpectedNumberOfBytesWritten
		}
	}
}

func (p *capitalizePipe) Finalize() {}

type pastTensePipe struct {
	out io.Writer
}

func ModifierPastTense(out io.Writer) Modifier {
	return &pastTensePipe{out: out}
}

func (p *pastTensePipe) Write(b []byte) (int, error) {
	return p.out.Write(b)
}

func (p *pastTensePipe) Finalize() {
	p.out.Write([]byte("ed"))
}
