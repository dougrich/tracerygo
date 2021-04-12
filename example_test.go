package tracerygo_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/dougrich/tracerygo"
)

func ExampleSingleStepInline() {
	g := tracerygo.RawGrammar{
		"origin":    []string{"hello #addressee#"},
		"addressee": []string{"world", "planet", "there"},
	}

	// name of field, index in field, seed
	r, err := g.Evaluate("origin", 0, 0)
	if err != nil {
		// this error happens if it couldn't parse the fields or is invalid in some way
		panic(err)
	}

	fmt.Println(r)
	// Output: hello world
}

func ExampleSingleStepJSON() {
	g := make(tracerygo.RawGrammar)

	if err := json.Unmarshal([]byte(`{"origin":["hello #addressee#"], "addressee":["world", "planet", "there"]}`), &g); err != nil {
		// an error occured parsing the json
		panic(err)
	}

	r, err := g.Evaluate("origin", 0, 0)
	if err != nil {
		// this error happens if it couldn't parse the fields or is invalid in some way
		panic(err)
	}

	fmt.Println(r)
	// Output: hello world
}

func ExampleMultiStep() {
	rawg := tracerygo.RawGrammar{
		"origin":    []string{"hello #addressee#"},
		"addressee": []string{"world", "planet", "there"},
	}

	g, err := tracerygo.Parse(rawg)
	if err != nil {
		// an error occured parsing the raw grammar
		panic(err)
	}

	r, err := g.Evaluate("origin", 0, 0)
	if err != nil {
		// an error occured evaluating the grammar
		panic(err)
	}

	fmt.Println(r)
	// Output: hello world
}

func ExampleCustomRandom() {
	rawg := tracerygo.RawGrammar{
		"origin":    []string{"hello #addressee#"},
		"addressee": []string{"world", "planet", "there"},
	}

	g, err := tracerygo.Parse(rawg)
	if err != nil {
		// an error occured parsing the raw grammar
		panic(err)
	}

	var sb strings.Builder
	e := tracerygo.NewEvaluation(&sb, tracerygo.WithRandom(rand.New(rand.NewSource(0))), tracerygo.WithGrammar(g))

	nodes, ok := g["origin"]
	if !ok || 0 >= len(nodes) {
		panic(errors.New("Missing origin"))
	}
	n := nodes[0]
	e.Evaluate(n)
	fmt.Println(sb.String())
	// Output: hello world
}

func ExampleCustomLookup() {
	rawg := tracerygo.RawGrammar{
		"origin": []string{"hello #addressee#"},
	}

	g, err := tracerygo.Parse(rawg)
	if err != nil {
		// an error occured parsing the raw grammar
		panic(err)
	}

	var sb strings.Builder
	e := tracerygo.NewEvaluation(
		&sb,
		tracerygo.WithGrammar(g),
		tracerygo.WithLookup(func(name string) (string, error) {
			if name == "addressee" {
				return "world", nil
			} else {
				return "", errors.New("Unknown lookup")
			}
		}),
	)

	nodes, ok := g["origin"]
	if !ok || 0 >= len(nodes) {
		panic(errors.New("Missing origin"))
	}
	n := nodes[0]
	e.Evaluate(n)
	fmt.Println(sb.String())
	// Output: hello world
}
