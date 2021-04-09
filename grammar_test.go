package tracerygo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	t.Run("string_ptr", func(t *testing.T) {
		result := stringNode("red").evaluate(nil)
		assert.Equal(t, "red", result)
	})
	t.Run("[]string", func(t *testing.T) {
		source := arrayNode{
			stringNode("red"),
			stringNode("blue"),
			stringNode("green"),
		}
		ctx := evaluationContext{
			next: func() uint64 { return 0 },
		}
		result := source.evaluate(&ctx)
		assert.Equal(t, "red", result)
		ctx = evaluationContext{
			next: func() uint64 { return math.MaxUint64 / 2 },
		}
		result = source.evaluate(&ctx)
		assert.Equal(t, "blue", result)
		ctx = evaluationContext{
			next: func() uint64 { return math.MaxUint64 },
		}
		result = source.evaluate(&ctx)
		assert.Equal(t, "green", result)
	})
	t.Run("format", func(t *testing.T) {
		f := formatNode{
			format:   "%s, %s, %s",
			children: []node{stringNode("red"), stringNode("blue"), stringNode("green")},
		}
		ctx := evaluationContext{
			next: func() uint64 { return 0 },
		}
		result := f.evaluate(&ctx)
		assert.Equal(t, "red, blue, green", result)
	})
	t.Run("modifier", func(t *testing.T) {
		f := modifierNode{
			root:     stringNode("duck"),
			modifier: func(r string) string { return "a " + r },
		}
		ctx := evaluationContext{
			next: func() uint64 { return 0 },
		}
		result := f.evaluate(&ctx)
		assert.Equal(t, "a duck", result)
	})
	t.Run("variable", func(t *testing.T) {
		v := variableNode{
			name: "alpha",
		}
		ctx := evaluationContext{
			next:      func() uint64 { return 0 },
			variables: map[string]string{"alpha": "duck"},
		}
		result := v.evaluate(&ctx)
		assert.Equal(t, "duck", result)
		v = variableNode{
			name: "beta",
		}
		result = v.evaluate(&ctx)
		assert.Equal(t, "#beta#", result)
	})
	t.Run("withVariable", func(t *testing.T) {
		v := withVariableNode{
			variableName: "alpha",
			variableNode: stringNode("duck"),
			child: variableNode{
				name: "alpha",
			},
		}
		ctx := evaluationContext{
			next: func() uint64 { return 0 },
		}
		result := v.evaluate(&ctx)
		assert.Equal(t, "duck", result)
	})
	t.Run("landscape", func(t *testing.T) {
		move := arrayNode{stringNode("spiral"), stringNode("twirl"), stringNode("curl"), stringNode("dance"), stringNode("twine"), stringNode("weave"), stringNode("meander"), stringNode("wander"), stringNode("flow")}
		path := arrayNode{stringNode("stream"), stringNode("brook"), stringNode("path"), stringNode("ravine"), stringNode("forest"), stringNode("fence"), stringNode("stone wall")}
		mood := arrayNode{stringNode("overcast"), stringNode("alight"), stringNode("clear"), stringNode("darkened"), stringNode("blue"), stringNode("shadowed"), stringNode("illuminated"), stringNode("silver"), stringNode("cool"), stringNode("warm"), stringNode("summer-warmed")}
		substance := arrayNode{stringNode("light"), stringNode("reflections"), stringNode("mist"), stringNode("shadow"), stringNode("darkness"), stringNode("brightness"), stringNode("gaiety"), stringNode("merriment")}
		nearby := arrayNode{
			formatNode{format: "beyond the %s", children: []node{variableNode{name: "path"}}},
			stringNode("far away"),
			stringNode("ahead"),
			stringNode("behind me"),
		}
		line := arrayNode{
			formatNode{
				format: "%s and %s, the %s was %s with %s",
				children: []node{
					modifierNode{
						root:     mood,
						modifier: modifierCapitalize,
					},
					mood,
					variableNode{name: "myPlace"},
					mood,
					substance,
				},
			},
			formatNode{
				format: "%s %s %s through the %s, filling me with %s",
				children: []node{
					modifierNode{
						root:     nearby,
						modifier: modifierCapitalize,
					},
					modifierNode{
						root:     variableNode{name: "myPlace"},
						modifier: modifierIndefiniteArticle,
					},
					modifierNode{
						root:     move,
						modifier: modifierPastTense,
					},
					path,
					substance,
				},
			},
		}
		origin := withVariableNode{
			variableName: "myPlace",
			variableNode: path,
			child:        line,
		}

		ctx := evaluationContext{
			next: func() uint64 { return 0 },
		}
		result := origin.evaluate(&ctx)
		assert.Equal(t, "Overcast and overcast, the stream was overcast with light", result)

		ctx = evaluationContext{
			next: func() uint64 { return math.MaxUint64 },
		}
		result = origin.evaluate(&ctx)
		assert.Equal(t, "Behind me a stone wall flowed through the stone wall, filling me with merriment", result)
	})
}
