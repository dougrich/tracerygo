package tracerygo

import (
	"errors"
)

var (
	ErrUnsupportedModifier = errors.New("Unsupported modifier")
)

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
		Parts: nil,
	}
	for _, v := range tokens {
		switch v.(type) {
		case string:
			n.Parts = append(n.Parts, v)
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
