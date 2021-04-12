package main

import (
	"encoding/json"
	"log"

	"github.com/dougrich/tracerygo"
)

func main() {
	rawg := make(tracerygo.RawGrammar)
	err := json.Unmarshal([]byte(`{"hexDigit":["0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F"],"digit":["0","1","2","3","4","5","6","7","8","9"],"color":["\\##hexDigit##hexDigit##hexDigit##hexDigit##hexDigit##hexDigit#"],"shape":["rect width=\"#num#\" height=\"#num#\" x=\"#num#\" y=\"#num#\"","circle cx=\"#num#\" cy=\"#num#\" r=\"#num#\"","polygon points=\"#num#,#num# #num#,#num# #num#,#num#\""],"num":["#digit##digit#"],"stroke":["stroke=\"#color#\" stroke-width=\"#digit#\" stroke-opacity=\"0.#digit#\"",""],"fill":["fill=\"#color#\" fill-opacity=\"0.#digit#\""],"makeShape":["<#shape# #stroke# #fill# />"],"bg":["<rect width='300' height='300' x='-100' y='-100' fill='#color#' />"],"origin":["<svg width=\"100\" height=\"100\">#bg##makeShape##makeShape##makeShape##makeShape##makeShape##makeShape##makeShape#</svg>"]}`), &rawg)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON %v", err)
	}

	result, err := rawg.Evaluate("origin", 0, 0)
	if err != nil {
		log.Fatalf("Error directly evaluating %v", err)
	}

	log.Print(result)
}
