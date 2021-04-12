package tracerygo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RawGrammar map[string][]string

func (rawg RawGrammar) Evaluate(name string, index int, seed int64) (string, error) {

	g, err := Parse(rawg)
	if err != nil {
		return "", err
	}

	return g.Evaluate("origin", 0, seed)
}

type tokenLookup struct {
	prefixes []variableDeclaration
	name     string
	suffixes []string
}

type variableDeclaration struct {
	name      string
	subtokens []interface{}
}

func toNode(tokens []interface{}) (Node, error) {
	n := Node{
		Variables: nil,
		Parts:     nil,
	}
	for _, v := range tokens {
		switch v.(type) {
		case string:
			n.Parts = append(n.Parts, v)
		case []variableDeclaration:
			v := v.([]variableDeclaration)
			for _, decl := range v {
				subn, err := toNode(decl.subtokens)
				if err != nil {
					return n, err
				}
				n.Variables = append(n.Variables, Variable{
					Key:   decl.name,
					Parts: subn.Parts,
				})
			}
		case tokenLookup:
			t := v.(tokenLookup)
			s := Substitution{
				Variables: nil,
				Modifiers: nil,
				Key:       t.name,
			}

			if len(t.prefixes) > 0 {
				s.Variables = make([]Variable, len(t.prefixes))
				for i, p := range t.prefixes {
					n, err := toNode(p.subtokens)
					if err != nil {
						return n, err
					}
					s.Variables[i] = Variable{
						Key:   p.name,
						Parts: n.Parts,
					}
				}
			}

			if len(t.suffixes) > 0 {
				s.Modifiers = make([]int, len(t.suffixes))
				for i, m := range t.suffixes {
					modifier, ok := modifierMap[m]
					if !ok {
						return n, UnsupportedModifierError{m}
					}
					s.Modifiers[i] = modifier
				}
			}

			n.Parts = append(n.Parts, s)
		}
	}
	return n, nil
}

func tokenize(input string) ([]interface{}, error) {
	var parts []interface{}
	inLookup := -1
	currentToken := ""
	var variableDeclarations []variableDeclaration

traversal:
	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '\\':
			if input[i+1] == '#' {
				currentToken += "#"
				i = i + 1
				continue traversal
			}
		case '[':
			// look ahead until we see the end, then break that string out to parse into a variable declaration
			for k := i + 1; k < len(input); k++ {
				switch input[k] {
				case ']':
					segments := strings.SplitN(input[i+1:k], ":", 2)
					name := segments[0]
					subparts, err := tokenize(segments[1])
					if err != nil {
						return nil, err
					}
					variableDeclarations = append(variableDeclarations, variableDeclaration{
						name,
						subparts,
					})
					i = k
					continue traversal
				}
			}
			return parts, UnmatchedSymbolError{i, "[", "]"}
		case '#':
			if inLookup >= 0 {
				inLookup = -1
				tokenParts := strings.Split(currentToken, ".")

				parts = append(parts, tokenLookup{variableDeclarations, tokenParts[0], tokenParts[1:]})
				variableDeclarations = nil
			} else {
				inLookup = i
				if variableDeclarations != nil {
					parts = append(parts, variableDeclarations)
				}
				if currentToken != "" {
					parts = append(parts, currentToken)
				}
			}
			currentToken = ""
			variableDeclarations = nil
			// this is here to explicity skip the default behavior
			continue traversal
		}
		currentToken += string(input[i])
	}

	if inLookup >= 0 {
		return nil, UnmatchedSymbolError{inLookup, "#", "#"}
	}

	if variableDeclarations != nil {
		parts = append(parts, variableDeclarations)
	}
	if currentToken != "" {
		parts = append(parts, currentToken)
	}

	return parts, nil
}

func (g *RawGrammar) UnmarshalJSON(data []byte) error {
	intermediate := make(map[string]interface{})
	local := make(RawGrammar)
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}

	for k, intermediateValue := range intermediate {
		switch v := intermediateValue.(type) {
		case string:
			local[k] = []string{v}
		case []interface{}:
			temparr := make([]string, len(v))
			for i, s := range v {
				switch subv := s.(type) {
				case string:
					temparr[i] = subv
				default:
					return FieldError{fmt.Sprintf("%s[%d]", k, i), ExpectationError{"a string", "something else"}}
				}
			}
			local[k] = temparr
		default:
			return FieldError{k, ExpectationError{"either an array of strings or a string", "something else"}}
		}
	}
	*g = local
	return nil
}

func Parse(g RawGrammar) (Grammar, error) {
	final := make(Grammar)
	for k, v := range g {
		var nodes []Node
		for i, raw := range v {
			tokens, err := tokenize(raw)
			if err != nil {
				return final, FieldError{fmt.Sprintf("%s[%d]", k, i), err}
			}
			node, err := toNode(tokens)
			if err != nil {
				return final, FieldError{fmt.Sprintf("%s[%d]", k, i), err}
			}
			nodes = append(nodes, node)
		}

		final[k] = nodes
	}
	return final, nil
}
