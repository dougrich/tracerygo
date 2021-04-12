# tracery-golang

This is an implementation of [Kate Compton](https://link.springer.com/chapter/10.1007/978-3-319-27036-4_14)'s tracery, an author-focused generative text tool. The reference implementation is writtin in JavaScript and found [here on github](https://github.com/galaxykate/tracery)

Tracery uses text-expansion: at a high level, it iterates over the string, replacing tokens with strings (that might also contain tokens) until no tokens remain. Additionally, at a high level, it is a pure function; relying only on the initial seed & grammar structure. You can imagine that a single call to it is tracing a ray through the possible probability space: one `trace` is one possible expansion of the grammar.

## Single Step Interface

```golang
func (rawg RawGrammar) Evaluate(name string, index int, seed int64) (string, error)
```

Some example use cases:
- generating a small example
- testing & experimenting with a small grammar set
- preparing the string for further manipulation

The single step abstracts away:
- what implementation of `*rand.Rand` to use
- the distinct steps of parsing & evaluating
- collecting the output

See ExampleSingleStepInline or ExampleSingleStepJSON.

## Multiple Step Interface

```golang
func Parse(g RawGrammar) (Grammar, error)
func (g Grammar) Evaluate(name string, index int, seed int64) (string, error)
func (g Grammar) StreamingEvaluate(out io.Writer, name string, index int, seed int64) error
```

Some example use cases:
- caching parsed values for repeated evaluations
- streaming large results

The multiple step abstracts away:
- what implementation of `*rand.Rand` to use

See ExampleMultiStep

## Full Interface

```golang
func NewEvaluation(out io.Writer, modifiers ...EvaluationModifier) (*Evaluation)
func WithRandom(rand *rand.Rand) EvaluationModifier
func WithGrammar(g Grammar) EvaluationModifier
```

Some example use cases:
- custom random function
- dynamically building up the grammar over time
- skipping parsing with inline values
- hooking into the variable lookup

See ExampleCustomRandom, ExampleRemoteLookup