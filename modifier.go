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
	modifierMap                       = map[string]int{
		"capitalize": ModifierCapitalizeIndex,
		"ed":         ModifierPastTenseIndex,
		"a":          ModifierIndefiniteArticleIndex,
	}
	ModifierCapitalizeIndex        = 1
	ModifierPastTenseIndex         = 2
	ModifierIndefiniteArticleIndex = 3
	modifierLookup                 = []ModifierFunc{
		nil,
		ModifierCapitalize,
		ModifierPastTense,
		ModifierIndefiniteArticle,
	}
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

type indefiniteArticlePipe struct {
	out  io.Writer
	done bool
}

func ModifierIndefiniteArticle(out io.Writer) Modifier {
	return &indefiniteArticlePipe{out: out, done: false}
}

func (p *indefiniteArticlePipe) Write(b []byte) (int, error) {
	if p.done {
		return p.out.Write(b)
	} else {
		p.done = true
		s := string(b)
		switch s[0] {
		case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
			s = "an " + s
		default:
			s = "a " + s
		}
		l, err := p.out.Write([]byte(s))
		if l == len(s) {
			return len(b), err
		} else {
			return 0, ErrUnexpectedNumberOfBytesWritten
		}
	}
}

func (p *indefiniteArticlePipe) Finalize() {}
