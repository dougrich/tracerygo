package tracerygo

import (
	"io"
	"math/rand"
	"strings"
)

type Node struct {
	Variables []Lookup
	Parts     []interface{}
}

type Lookup struct {
	Key   string
	Value string
}

type Substitution struct {
	Variables []Lookup
	Modifiers []ModifierFunc
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
		Lookup: make(map[string][]Node),
	}
	for _, m := range modifiers {
		m(e)
	}
	if e.rand == nil {
		e.rand = rand.New(rand.NewSource(0))
	}
	return e, nil
}

func WithRandom(rand *rand.Rand) EvaluationModifier {
	return func(e *Evaluation) {
		e.rand = rand
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
			varn := sube.EvaluateLookup(v.Value)
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
					modifiers[i] = m(pipe)
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
