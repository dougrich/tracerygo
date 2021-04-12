package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/dougrich/tracerygo"
)

func main() {
	rawg := make(tracerygo.RawGrammar)
	bytes, err := ioutil.ReadFile("grammar.json")
	if err != nil {
		log.Fatalf("Error reading file %v", err)
	}
	err = json.Unmarshal(bytes, &rawg)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON %v", err)
	}

	result, err := rawg.Evaluate("origin", 0, 0)
	if err != nil {
		log.Fatalf("Error directly evaluating %v", err)
	}

	log.Print(result)
}
