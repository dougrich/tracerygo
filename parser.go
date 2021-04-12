package tracerygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	ErrUnsupportedModifier = errors.New("Unsupported modifier")
)

type RawGrammar map[string][]string

func (rawg RawGrammar) Evaluate(name string, index int) (string, error) {

	g, err := Parse(rawg)
	if err != nil {
		return "", err
	}

	return g.SEvaluate("origin", 0, time.Now().UnixNano())
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
						return n, fmt.Errorf("in %s, unsupported modifier %s", "", m)
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
			log.Printf("TRACE: lookahead started at %d", i)
			for k := i + 1; k < len(input); k++ {
				switch input[k] {
				case ']':
					log.Printf("TRACE: lookahead ended at %d, captured %s", i, input[i+1:k])
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
			return parts, fmt.Errorf("Unmatched '[', started at %d", i)
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
			// this is here to explicity skip the default behavior
			continue traversal
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
					return fmt.Errorf("Error in field '%s': expected an array of strings, but element %d doesn't appear to be a string", k, i)
				}
			}
			local[k] = temparr
		default:
			return fmt.Errorf("Error in field '%s': expected either a string or an array of strings", k)
		}
	}
	*g = local
	return nil
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
