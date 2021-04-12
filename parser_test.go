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
		tokenLookup{[]variableDeclaration{{"myworld", []interface{}{"cool ", tokenLookup{nil, "world", []string{}}}}}, "myworld", nil},
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

	n, err = toNode([]interface{}{
		[]variableDeclaration{
			{
				"neat",
				[]interface{}{
					tokenLookup{nil, "test", nil},
				},
			},
		},
		tokenLookup{nil, "world", []string{}},
	})
	assert.Nil(err)
	assert.Equal(n, Node{
		Variables: []Lookup{
			{
				"neat",
				[]interface{}{
					Substitution{
						Variables: nil,
						Modifiers: nil,
						Key:       "test",
					},
				},
			},
		},
		Parts: []interface{}{
			Substitution{
				Variables: nil,
				Modifiers: nil,
				Key:       "world",
			},
		},
	})
}

func TestTokenize(t *testing.T) {
	assert := assert.New(t)
	var parts []interface{}
	var err error

	parts, err = tokenize("red")
	if assert.Nil(err) {
		assert.Equal([]interface{}{"red"}, parts)
	}

	parts, err = tokenize("red #type#")
	if assert.Nil(err) {
		assert.Equal([]interface{}{"red ", tokenLookup{nil, "type", []string{}}}, parts)
	}

	parts, err = tokenize("red \\##type#")
	if assert.Nil(err) {
		assert.Equal([]interface{}{"red #", tokenLookup{nil, "type", []string{}}}, parts)
	}

	parts, err = tokenize("red #type.a.b#")
	if assert.Nil(err) {
		assert.Equal([]interface{}{"red ", tokenLookup{nil, "type", []string{"a", "b"}}}, parts)
	}

	parts, err = tokenize("red #[myVar:#neat#]type#")
	if assert.Nil(err) {
		assert.Equal(
			[]interface{}{
				"red ",
				tokenLookup{
					[]variableDeclaration{
						{
							"myVar",
							[]interface{}{tokenLookup{nil, "neat", []string{}}},
						},
					},
					"type",
					[]string{},
				},
			},
			parts,
		)
	}

	parts, err = tokenize("[myVar:#neat#]#type#")
	if assert.Nil(err) {
		assert.Equal(
			[]interface{}{
				[]variableDeclaration{
					{
						"myVar",
						[]interface{}{tokenLookup{nil, "neat", []string{}}},
					},
				},
				tokenLookup{
					nil,
					"type",
					[]string{},
				},
			},
			parts,
		)
	}

	// this is an extract from the sci-fi example that caused an issue parsing
	parts, err = tokenize("#[mcArt:#artForm#][mcBoss:#boss#]artPlot#")
	if assert.Nil(err) {
		assert.Equal(
			[]interface{}{
				tokenLookup{
					[]variableDeclaration{
						{
							"mcArt",
							[]interface{}{tokenLookup{nil, "artForm", []string{}}},
						},
						{
							"mcBoss",
							[]interface{}{tokenLookup{nil, "boss", []string{}}},
						},
					},
					"artPlot",
					[]string{},
				},
			},
			parts,
		)
	}
}
