package tracerygo

import (
	"errors"
	"io"
	"math/rand"
	"strings"
)

// This represents a node in the evaluation tree
type Node struct {
	// An array of variable declarations that apply to this node and it's children
	Variables []Variable
	// The parts to be evaluated; this is untyped but can contain strings, Substitutions, Evaluations, etc.
	Parts []interface{}
}

// This represents a variable definition.
type Variable struct {
	// This is the key that can be used later (i.e. `Variable{"myVar", []interface{}{"value"}}` is equivalent to "[myVar:value]")
	Key string
	// The parts to be evaluated when it is looked up; this could be a lookup, a string, or similar. It's evaluated like a Node
	Parts []interface{}
}

// This represents a substitution, e.g. '#myVar#' or '#[myVar:#sub#]value.s'
type Substitution struct {
	// An array of variable declarations that apply to this lookup and it's children
	Variables []Variable
	// An array of modifier indexes; ideally this should be functions, but for testing purposes we use indexes
	Modifiers []int
	// The key to lookup and replace this substitution with
	Key string
}

// This is the bundled context for a single evaluation.
type Evaluation struct {
	// The grammar defined alongside the current evaluation that might be drilled into
	Grammar map[string][]Node

	// note that these variables are private; this is done mostly to encourage going through the 'WithRandom' interface and to prevent long-lived Evaluations from kicking around
	// this is a custom lookup function
	lookup LookupFunction
	// this is a custom random intreface
	rand *rand.Rand
	// this is the output stream
	out io.Writer
}

// An evaluation modifier, when passed in to create the evaluation, modifies it's internal state on creation. This can be used to give an optional paramter or some configuration value
type EvaluationModifier func(*Evaluation)

// This creates a new evaluation context for a specific stream. The modifiers are run in order and sane defaults are provided for optional values
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

// This provides a custom random interface to an evaluation context
func WithRandom(rand *rand.Rand) EvaluationModifier {
	return func(e *Evaluation) {
		e.rand = rand
	}
}

// This provides a custom grammar to an evaluation context
func WithGrammar(g Grammar) EvaluationModifier {
	return func(e *Evaluation) {
		e.Grammar = g
	}
}

// This is a lookup function; it takes a string that is not found in the existing grammar and returns a string that should be used, or an error if it can't be found
type LookupFunction func(string) (string, error)

// This provides a custom lookup function to an evaluation context
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

// This evaluates an entire node, writing it to the underlying stream directly
func (e *Evaluation) Evaluate(n Node) error {
	if len(n.Variables) != 0 {
		e = e.clone(nil, n.Variables)
	}
	for _, abstract := range n.Parts {
		switch v := abstract.(type) {
		case string:
			if _, err := e.out.Write([]byte(v)); err != nil {
				return ErrorStreamWrite{err}
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

// This evaluates a specific name as if it were looking it up, writing it to the underlying stream directly
func (e *Evaluation) EvaluateName(name string) (Node, error) {
	nodes, ok := e.Grammar[name]
	if !ok || len(nodes) == 0 {
		// not found; do we have a lookup function?
		if e.lookup != nil {
			value, err := e.lookup(name)
			if err != nil {
				return Node{}, ErrorLookup{name, err}
			}
			return Node{Parts: []interface{}{value}}, nil
		} else {
			return Node{}, ErrorNameNotFound{name}
		}
	}
	i := e.rand.Intn(len(nodes))
	return nodes[i], nil
}

// This represents a parsed and ready to use grammar
type Grammar map[string][]Node

// This evaluates and directly streams it out to a specified writer
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

// This calls StreamingEvaluate under the hood and buffers it to a string before returning
func (g Grammar) Evaluate(name string, index int, seed int64) (string, error) {
	var sb strings.Builder
	err := g.StreamingEvaluate(&sb, name, index, seed)
	return sb.String(), err
}
