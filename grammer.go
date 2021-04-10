package tracerygo

import (
	"io"
	"math/rand"
	"strings"
)

type Node struct {
	Parts []interface{}
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

func (e *Evaluation) clone(out io.Writer, lookups map[string]Node) *Evaluation {
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

	if lookups != nil {
		sube.Lookup = make(map[string][]Node)
		for k, v := range e.Lookup {
			sube.Lookup[k] = v
		}
		for k, v := range lookups {
			sube.Lookup[k] = []Node{v}
		}
	}

	return sube
}

func (e *Evaluation) Evaluate(n Node) {
	for _, abstract := range n.Parts {
		switch v := abstract.(type) {
		case string:
			e.out.Write([]byte(v))
		case Substitution:
			var pipe io.Writer
			var modifiers []Modifier
			var variables map[string]Node
			n := e.EvaluateLookup(v.Key)

			if len(v.Modifiers) != 0 {
				modifiers = make([]Modifier, len(v.Modifiers))
				pipe = e.out
				for i, m := range v.Modifiers {
					modifiers[i] = m(pipe)
					pipe = modifiers[i]
				}
			}

			if len(v.Variables) != 0 {
				variables = make(map[string]Node)
				for _, variable := range v.Variables {
					// fully evaluate the variable now and cache the result
					varn := e.EvaluateLookup(variable.Value)
					var sb strings.Builder
					sube := e.clone(&sb, variables)
					sube.Evaluate(varn)
					variables[variable.Key] = Node{Parts: []interface{}{sb.String()}}
				}
			}

			sube := e.clone(pipe, variables)
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
