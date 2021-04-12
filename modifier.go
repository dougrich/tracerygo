package tracerygo

import (
	"io"
	"strings"
)

// This is a 'in between' function which holds onto the state needed for modifiers ('.ed' and similar)
type ModifierFunc func(io.Writer) Modifier

//This is the interface between the in-between state; the 'Finalize' method enables suffixes
type Modifier interface {
	io.Writer
	Finalize() error
}

var (
	modifierMap = map[string]int{
		"capitalize": modifierCapitalizeIndex,
		"ed":         modifierPastTenseIndex,
		"a":          modifierIndefiniteArticleIndex,
		"s":          modifierPluralizeIndex,
	}
	modifierCapitalizeIndex        = 1
	modifierPastTenseIndex         = 2
	modifierIndefiniteArticleIndex = 3
	modifierPluralizeIndex         = 4
	modifierLookup                 = []ModifierFunc{
		nil,
		ModifierCapitalize,
		ModifierPastTense,
		ModifierIndefiniteArticle,
		ModifierPluralize,
	}
)

type capitalizePipe struct {
	out  io.Writer
	done bool
}

// This returns a modifier for capitalizing
func ModifierCapitalize(out io.Writer) Modifier {
	return &capitalizePipe{out: out, done: false}
}

// This writes to the underlying stream; the first set of bytes that it gets it will try to capitalize the first letter of
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

// This matches the interface but does nothing
func (p *capitalizePipe) Finalize() error {
	return nil
}

type pastTensePipe struct {
	out io.Writer
}

// This returns a modifier for turning a verb into the past tense
func ModifierPastTense(out io.Writer) Modifier {
	return &pastTensePipe{out: out}
}

// This is just a passthrough
func (p *pastTensePipe) Write(b []byte) (int, error) {
	return p.out.Write(b)
}

// This appends an 'ed' top the stream
func (p *pastTensePipe) Finalize() error {
	_, err := p.out.Write([]byte("ed"))
	return err
}

type indefiniteArticlePipe struct {
	out  io.Writer
	done bool
}

// This returns a modifier for prefixing the indefinite article to a noun
func ModifierIndefiniteArticle(out io.Writer) Modifier {
	return &indefiniteArticlePipe{out: out, done: false}
}

// This writes to the underlying stream; the first set of bytes are inspected and the suitable 'a' or 'an' is placed in front
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

// This matches the interface but does nothing
func (p *indefiniteArticlePipe) Finalize() error {
	return nil
}

type pluralizePipe struct {
	out io.Writer
}

// This returns a modifier for pluralizing a noun
func ModifierPluralize(out io.Writer) Modifier {
	return &pluralizePipe{out: out}
}

// This is just a passthrough
func (p *pluralizePipe) Write(b []byte) (int, error) {
	return p.out.Write(b)
}

// This appends an 's' to the stream
func (p *pluralizePipe) Finalize() error {
	_, err := p.out.Write([]byte("s"))
	return err
}
