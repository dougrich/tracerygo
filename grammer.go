package tracerygo

import (
	"fmt"
	"math"
	"strings"
)

var (
	modifierCapitalize modifierFunc = func(s string) string {
		return strings.ToUpper(string(s[0])) + s[1:]
	}
	modifierIndefiniteArticle modifierFunc = func(s string) string {
		switch s[0] {
		case 'a', 'e', 'o', 'u', 'i', 'A', 'E', 'O', 'U', 'I':
			return "an " + s
		default:
			return "a " + s
		}
	}
	modifierPastTense modifierFunc = func(s string) string {
		return s + "ed"
	}
)

type node interface {
	evaluate(c *evaluationContext) string
}

type stringNode string

func (s stringNode) evaluate(_ *evaluationContext) string {
	s2 := s
	return string(s2)
}

type arrayNode []node

func (s arrayNode) evaluate(ctx *evaluationContext) string {
	i := ctx.choose(len(s))
	return s[i].evaluate(ctx)
}

type formatNode struct {
	format   string
	children []node
}

func (f formatNode) evaluate(ctx *evaluationContext) string {
	margs := make([]interface{}, len(f.children))
	for i, c := range f.children {
		margs[i] = c.evaluate(ctx)
	}
	return fmt.Sprintf(f.format, margs...)
}

type modifierFunc (func(string) string)
type modifierNode struct {
	root     node
	modifier modifierFunc
}

func (m modifierNode) evaluate(ctx *evaluationContext) string {
	root := m.root.evaluate(ctx)
	return m.modifier(root)
}

type variableNode struct {
	name string
}

func (v variableNode) evaluate(ctx *evaluationContext) string {
	value, ok := ctx.variables[v.name]
	if !ok {
		return "#" + v.name + "#"
	} else {
		return value
	}
}

type withVariableNode struct {
	variableName string
	variableNode node
	child        node
}

func (w withVariableNode) evaluate(ctx *evaluationContext) string {
	ctx2 := ctx.clone()
	ctx2.variables[w.variableName] = w.variableNode.evaluate(ctx)
	return w.child.evaluate(ctx2)
}

type nextFunc (func() uint64)

type evaluationContext struct {
	next      nextFunc
	variables map[string]string
}

func (e *evaluationContext) clone() *evaluationContext {
	e2 := evaluationContext{
		next:      e.next,
		variables: make(map[string]string),
	}
	for k, v := range e.variables {
		e2.variables[k] = v
	}
	return &e2
}

func (e *evaluationContext) choose(len int) int {
	n := e.next()
	return int(float64(len) * 0.99999999999 * float64(n) / (float64(math.MaxUint64)))
}
