package tracerygo

import (
	"fmt"
	"math"
	"strings"
	"errors"
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

type variableLocation struct {
	start int
	end int
}

func parse(input string) (node, error) {
	if input[0] == '[' {
		// scan for the end input
		next := strings.IndexOf(input, "]")
		if next == -1 {
			return nil, errors.New("String substitution started; but missing end")
		}
		variableDeclarationPhrase := input[1:next]
		variableDeclarations := strings.Split(variableDeclarationPhrase, ",")
		n, err := parse(input[next+1:])
		if err != nil {
			return err
		}
		for _, v := variableDeclarations {
			
		}
	}
	var variables []variableLocation
	startSubstitution := -1

	for i, r := range input {
		switch r {
		case '#':
			if startSubstitution == -1 {
				startSubstitution = i		
			} else {
				variables = append(variables, variableLocation{ startSubstitution, i })
				startSubstitution = -1
			}
		}
	}

	if len(variables) == 0 {
		return stringNode(input), nil
	}

	last := 0
	format := ""
	children := make([]node, len(variables))
	var errs []error
	for i, v := range variables {
		format += input[last:v.start]
		format += "%s"
		last = v.end + 1
		name := input[v.start+1:v.end]
		mods := strings.Split(name, ".")

		var n node = variableNode{ mods[0] }
		for _, key := range mods[1:] {
			switch key {
			case "a":
				n = modifierNode{ root: n, modifier: modifierIndefiniteArticle }
			case "capitalize":
				n = modifierNode{ root: n, modifier: modifierCapitalize }
			case "ed":
				n = modifierNode{ root: n, modifier: modifierPastTense }
			default:
				errs = append(errs, errors.New("Unrecognized modifier"))
			}
		}
		children[i] = n
		
	}
	if errs != nil {
		return nil, errors.New("Errors occured!")
	}
	format += input[last:]
	return formatNode{
		format,
		children,
	}, nil
}

func parseArray(arr []string) (node, error) {
	var n arrayNode
	for _, s := range arr {
		nc, err := parse(s)
		if err != nil {
			return nil, err
		}
		n = append(n, nc)
	}
	return n, nil
}

type grammar map[string][]string

func parseGrammar(g grammar) (node, error) {
	n, err := parse(g["origin"][0])
	if err != nil {
		return nil, err
	}

	for k, v := range g {
		if k == "origin" {
			// skip over alternate origins
			continue
		}

		a, err := parseArray(v)
		if err != nil {
			return nil, err
		}

		n = withVariableNode{
			variableName: k,
			variableNode: a,
			child: n,
		}
	}

	return n, nil
}