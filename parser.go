package tracerygo

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedModifier = errors.New("Unsupported modifier")
)

type RawGrammar map[string][]string

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
				n.Variables = append(n.Variables, Lookup{
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
				s.Variables = make([]Lookup, len(t.prefixes))
				for i, p := range t.prefixes {
					n, err := toNode(p.subtokens)
					if err != nil {
						return n, err
					}
					s.Variables[i] = Lookup{
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
						return n, ErrUnsupportedModifier
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
	inLookup := false
	currentToken := ""
	var variableDeclarations []variableDeclaration

	for i := 0; i < len(input); i++ {
		switch input[i] {
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
				}
			}
			continue
		case '#':
			if inLookup {
				inLookup = false
				tokenParts := strings.Split(currentToken, ".")

				parts = append(parts, tokenLookup{variableDeclarations, tokenParts[0], tokenParts[1:]})
				variableDeclarations = nil
			} else {
				inLookup = true
				if variableDeclarations != nil {
					parts = append(parts, variableDeclarations)
				}
				if currentToken != "" {
					parts = append(parts, currentToken)
				}
			}
			currentToken = ""
			variableDeclarations = nil
			continue
		}
		currentToken += string(input[i])
	}

	if variableDeclarations != nil {
		parts = append(parts, variableDeclarations)
	}
	if currentToken != "" {
		parts = append(parts, currentToken)
	}

	return parts, nil
}

func Parse(g RawGrammar) (Grammar, error) {
	final := make(Grammar)
	for k, v := range g {
		var nodes []Node
		for _, raw := range v {
			tokens, err := tokenize(raw)
			if err != nil {
				return final, err
			}
			node, err := toNode(tokens)
			if err != nil {
				return final, err
			}
			nodes = append(nodes, node)
		}

		final[k] = nodes
	}
	return final, nil
}
