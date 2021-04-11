package tracerygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToNode(t *testing.T) {
	assert := assert.New(t)

	n, err := toNode([]interface{}{
		"hello ",
		"world",
	})
	assert.Nil(err)
	assert.Equal(n, Node{
		Parts: []interface{}{
			"hello ",
			"world",
		},
	})

	n, err = toNode([]interface{}{
		"hello ",
		tokenLookup{nil, "world", nil},
	})
	assert.Nil(err)
	assert.Equal(n, Node{
		Parts: []interface{}{
			"hello ",
			Substitution{
				Variables: nil,
				Modifiers: nil,
				Key:       "world",
			},
		},
	})

	n, err = toNode([]interface{}{
		"hello ",
		tokenLookup{[]variableDeclaration{{"myworld", []interface{}{"cool ", tokenLookup{nil, "world", nil}}}}, "myworld", nil},
	})
	assert.Nil(err)
	assert.Equal(n, Node{
		Parts: []interface{}{
			"hello ",
			Substitution{
				Variables: []Lookup{
					{
						Key: "myworld",
						Parts: []interface{}{
							"cool ",
							Substitution{Key: "world"},
						},
					},
				},
				Modifiers: nil,
				Key:       "myworld",
			},
		},
	})

	n, err = toNode([]interface{}{
		"hello ",
		tokenLookup{nil, "world", []string{"capitalize"}},
	})
	assert.Nil(err)
	assert.Equal(n, Node{
		Parts: []interface{}{
			"hello ",
			Substitution{
				Variables: nil,
				Modifiers: []int{
					ModifierCapitalizeIndex,
				},
				Key: "world",
			},
		},
	})
}
