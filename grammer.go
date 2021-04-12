package tracerygo

import (
	"errors"
	"io"
	"math/rand"
	"strings"
)

type Node struct {
	Variables []Variable
	Parts     []interface{}
}

type Variable struct {
	Key   string
	Parts []interface{}
}

type Substitution struct {
	Variables []Variable
	Modifiers []int
	Key       string
}

type Evaluation struct {
	Grammar map[string][]Node
	lookup  LookupFunction
	rand    *rand.Rand
	out     io.Writer
}

type EvaluationModifier func(*Evaluation)

func NewEvaluation(out io.Writer, modifiers ...EvaluationModifier) *Evaluation {
	e := &Evaluation{
		Grammar: nil,
		lookup:  nil,
		out:     out,
		rand:    nil,
	}
	for _, m := range modifiers {
		m(e)
	}
	if e.rand == nil {
		e.rand = rand.New(rand.NewSource(0))
	}
	if e.Grammar == nil {
		e.Grammar = make(map[string][]Node)
	}
	return e
}

func WithRandom(rand *rand.Rand) EvaluationModifier {
	return func(e *Evaluation) {
		e.rand = rand
	}
}

func WithGrammar(g Grammar) EvaluationModifier {
	return func(e *Evaluation) {
		e.Grammar = g
	}
}

type LookupFunction func(string) (string, error)

func WithLookup(fn LookupFunction) EvaluationModifier {
	return func(e *Evaluation) {
		e.lookup = fn
	}
}

func (e *Evaluation) clone(out io.Writer, lookups []Variable) *Evaluation {
	// shortcut if we're cloning but don't actually make any changes
	if out == nil && lookups == nil {
		return e
	}

	sube := &Evaluation{
		out:     e.out,
		rand:    e.rand,
		Grammar: e.Grammar,
	}

	if out != nil {
		sube.out = out
	}

	if len(lookups) != 0 {
		sube.Grammar = make(map[string][]Node)
		for k, v := range e.Grammar {
			sube.Grammar[k] = v
		}
		for _, v := range lookups {
			varn := Node{Parts: v.Parts}
			var sb strings.Builder
			t := sube.out
			sube.out = &sb
			sube.Evaluate(varn)
			sube.Grammar[v.Key] = []Node{{Parts: []interface{}{sb.String()}}}
			sube.out = t
		}
	}

	return sube
}

func (e *Evaluation) Evaluate(n Node) error {
	if len(n.Variables) != 0 {
		e = e.clone(nil, n.Variables)
	}
	for _, abstract := range n.Parts {
		switch v := abstract.(type) {
		case string:
			if _, err := e.out.Write([]byte(v)); err != nil {
				return StreamWriteError{err}
			}
		case Substitution:
			var pipe io.Writer
			var modifiers []Modifier
			n, err := e.EvaluateName(v.Key)
			if err != nil {
				return err
			}

			if len(v.Modifiers) != 0 {
				modifiers = make([]Modifier, len(v.Modifiers))
				pipe = e.out
				for i, m := range v.Modifiers {
					modifiers[i] = modifierLookup[m](pipe)
					pipe = modifiers[i]
				}
			}

			sube := e.clone(pipe, v.Variables)
			if err := sube.Evaluate(n); err != nil {
				return err
			}

			for _, m := range modifiers {
				if err := m.Finalize(); err != nil {
					return err
				}
			}
		default:
		}
	}
	return nil
}

func (e *Evaluation) EvaluateName(name string) (Node, error) {
	nodes, ok := e.Grammar[name]
	if !ok || len(nodes) == 0 {
		// not found; do we have a lookup function?
		if e.lookup != nil {
			value, err := e.lookup(name)
			if err != nil {
				return Node{}, LookupError{name, err}
			}
			return Node{Parts: []interface{}{value}}, nil
		} else {
			return Node{}, NameNotFoundError{name}
		}
	}
	i := e.rand.Intn(len(nodes))
	return nodes[i], nil
}

type Grammar map[string][]Node

func (g Grammar) StreamingEvaluate(out io.Writer, name string, index int, seed int64) error {
	e := NewEvaluation(out, WithRandom(rand.New(rand.NewSource(seed))), WithGrammar(g))

	nodes, ok := g[name]
	if !ok {
		return errors.New("Key not found")
	}
	if index < 0 || index >= len(nodes) {
		return errors.New("Index out of bounds")
	}
	n := nodes[index]
	return e.Evaluate(n)
}

func (g Grammar) Evaluate(name string, index int, seed int64) (string, error) {
	var sb strings.Builder
	err := g.StreamingEvaluate(&sb, name, index, seed)
	return sb.String(), err
}
