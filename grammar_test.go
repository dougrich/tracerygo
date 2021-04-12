package tracerygo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		result := Node{
			Parts: []interface{}{
				"hello ",
				"world",
			},
		}
		var sb strings.Builder
		ctx := NewEvaluation(&sb)
		ctx.Evaluate(result)
		assert.Equal(t, "hello world", sb.String())
	})
	t.Run("substitution", func(t *testing.T) {
		result := Node{
			Parts: []interface{}{
				"hello ",
				Substitution{
					Key: "world",
				},
			},
		}
		var sb strings.Builder
		ctx := NewEvaluation(&sb)
		ctx.Grammar["world"] = []Node{
			{
				Parts: []interface{}{
					"world",
				},
			},
		}
		ctx.Evaluate(result)
		assert.Equal(t, "hello world", sb.String())
	})
	t.Run("substitution with prefix odifier", func(t *testing.T) {
		result := Node{
			Parts: []interface{}{
				"hello ",
				Substitution{
					Modifiers: []int{
						ModifierCapitalizeIndex,
					},
					Key: "world",
				},
			},
		}
		var sb strings.Builder
		ctx := NewEvaluation(&sb)
		ctx.Grammar["world"] = []Node{
			{
				Parts: []interface{}{
					"world",
				},
			},
		}
		ctx.Evaluate(result)
		assert.Equal(t, "hello World", sb.String())
	})
	t.Run("substitution with suffix modifier", func(t *testing.T) {
		result := Node{
			Parts: []interface{}{
				"hello ",
				Substitution{
					Modifiers: []int{
						ModifierPastTenseIndex,
					},
					Key: "world",
				},
			},
		}
		var sb strings.Builder
		ctx := NewEvaluation(&sb)
		ctx.Grammar["world"] = []Node{
			{
				Parts: []interface{}{
					"world",
				},
			},
		}
		ctx.Evaluate(result)
		assert.Equal(t, "hello worlded", sb.String())
	})
	t.Run("substitution with variables", func(t *testing.T) {
		result := Node{
			Parts: []interface{}{
				"hello ",
				Substitution{
					Variables: []Variable{
						{"myWorld", []interface{}{Substitution{Key: "world2"}}},
					},
					Key: "world",
				},
			},
		}
		var sb strings.Builder
		ctx := NewEvaluation(&sb)
		ctx.Grammar["world"] = []Node{
			{
				Parts: []interface{}{
					Substitution{
						Key: "myWorld",
					},
					" and ",
					Substitution{
						Key: "myWorld",
					},
				},
			},
		}
		ctx.Grammar["world2"] = []Node{
			{
				Parts: []interface{}{
					"world",
				},
			},
		}
		ctx.Evaluate(result)
		assert.Equal(t, "hello world and world", sb.String())
	})
	t.Run("substitution with global variables", func(t *testing.T) {
		result := Node{
			Variables: []Variable{
				{"myWorld", []interface{}{Substitution{Key: "world2"}}},
			},
			Parts: []interface{}{
				"hello ",
				Substitution{
					Key: "world",
				},
				", and ",
				Substitution{
					Key: "world",
				},
			},
		}
		var sb strings.Builder
		ctx := NewEvaluation(&sb)
		ctx.Grammar["world"] = []Node{
			{
				Parts: []interface{}{
					Substitution{
						Key: "myWorld",
					},
					" and ",
					Substitution{
						Key: "myWorld",
					},
				},
			},
		}
		ctx.Grammar["world2"] = []Node{
			{
				Parts: []interface{}{
					"world",
				},
			},
		}
		ctx.Evaluate(result)
		assert.Equal(t, "hello world and world, and world and world", sb.String())
	})
}
