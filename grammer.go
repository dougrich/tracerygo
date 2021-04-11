package tracerygo

import (
	"io"
	"math/rand"
	"strings"
	"errors"
)

type Node struct {
	Variables []Lookup
	Parts     []interface{}
}

type Lookup struct {
	Key   string
	Parts []interface{}
}

type Substitution struct {
	Variables []Lookup
	Modifiers []int
	Key       string
}

type Evaluation struct {
	Lookup map[string][]Node
	rand   *rand.Rand
	out    io.Writer
}

type EvaluationModifier func(*Evaluation)

func NewEvaluation(out io.Writer, modifiers ...EvaluationModifier) (*Evaluation, error) {
	e := &Evaluation{
		out:    out,
		rand:   nil,
		Lookup: nil,
	}
	for _, m := range modifiers {
		m(e)
	}
	if e.rand == nil {
		e.rand = rand.New(rand.NewSource(0))
	}
	if e.Lookup == nil {
		e.Lookup = make(map[string][]Node)
	}
	return e, nil
}

func WithRandom(rand *rand.Rand) EvaluationModifier {
	return func(e *Evaluation) {
		e.rand = rand
	}
}

func WithGrammar(g Grammar) EvaluationModifier {
	return func(e *Evaluation) {
		e.Lookup = g
	}
}

func (e *Evaluation) clone(out io.Writer, lookups []Lookup) *Evaluation {
	// shortcut if we're cloning but don't actually make any changes
	if out == nil && lookups == nil {
		return e
	}

	sube := &Evaluation{
		out:    e.out,
		rand:   e.rand,
		Lookup: e.Lookup,
	}

	if out != nil {
		sube.out = out
	}

	if len(lookups) != 0 {
		sube.Lookup = make(map[string][]Node)
		for k, v := range e.Lookup {
			sube.Lookup[k] = v
		}
		for _, v := range lookups {
			varn := Node{Parts: v.Parts}
			var sb strings.Builder
			t := sube.out
			sube.out = &sb
			sube.Evaluate(varn)
			sube.Lookup[v.Key] = []Node{{Parts: []interface{}{sb.String()}}}
			sube.out = t
		}
	}

	return sube
}

func (e *Evaluation) Evaluate(n Node) {
	if len(n.Variables) != 0 {
		e = e.clone(nil, n.Variables)
	}
	for _, abstract := range n.Parts {
		switch v := abstract.(type) {
		case string:
			e.out.Write([]byte(v))
		case Substitution:
			var pipe io.Writer
			var modifiers []Modifier
			n := e.EvaluateLookup(v.Key)

			if len(v.Modifiers) != 0 {
				modifiers = make([]Modifier, len(v.Modifiers))
				pipe = e.out
				for i, m := range v.Modifiers {
					modifiers[i] = modifierLookup[m](pipe)
					pipe = modifiers[i]
				}
			}

			sube := e.clone(pipe, v.Variables)
			sube.Evaluate(n)

			for _, m := range modifiers {
				m.Finalize()
			}
		default:
		}
	}
}

func (e *Evaluation) EvaluateLookup(name string) Node {
	nodes, ok := e.Lookup[name]
	if !ok || len(nodes) == 0 {
		return Node{Parts: []interface{}{"#" + name + "#"}}
	}
	i := e.rand.Intn(len(nodes))
	return nodes[i]
}

type Grammar map[string][]Node

func (g Grammar) FEvaluate(out io.Writer, name string, index int, seed int64) (error) {
	e, err := NewEvaluation(out, WithRandom(rand.New(rand.NewSource(seed))), WithGrammar(g))
	if err != nil {
		return err
	}
	nodes, ok := g[name]
	if !ok {
		return errors.New("Key not found")
	}
	if index < 0 || index >= len(nodes) {
		return errors.New("Index out of bounds")
	}
	n := nodes[index]
	e.Evaluate(n)
	return nil
}

func (g Grammar) SEvaluate(name string, index int, seed int64) (string, error) {
	var sb strings.Builder
	err := g.FEvaluate(&sb, name, index, seed)
	return sb.String(), err
}
