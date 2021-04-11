package main

import (
	"log"
	"time"

	"github.com/dougrich/tracerygo"
)

func main() {
	rawg := tracerygo.RawGrammar{
		"origin":    []string{"[myPlace:#path#]#line#"},
		"line":      []string{"#mood.capitalize# and #mood#, the #myPlace# was #mood# with #substance#", "#nearby.capitalize# #myPlace.a# #move.ed# through the #path#, filling me with #substance#"},
		"nearby":    []string{"beyond the #path#", "far away", "ahead", "behind me"},
		"substance": []string{"light", "reflections", "mist", "shadow", "darkness", "brightness", "gaiety", "merriment"},
		"mood":      []string{"overcast", "alight", "clear", "darkened", "blue", "shadowed", "illuminated", "silver", "cool", "warm", "summer-warmed"},
		"path":      []string{"stream", "brook", "path", "ravine", "forest", "fence", "stone wall"},
		"move":      []string{"spiral", "twirl", "curl", "dance", "twine", "weave", "meander", "wander", "flow"},
	}

	g, err := tracerygo.Parse(rawg)
	if err != nil {
		log.Fatalf("Error parsing %v", err)
	}

	result, err := g.SEvaluate("origin", 0, time.Now().UnixNano())
	if err != nil {
		log.Fatalf("Error evaluation %v", err)
	}

	log.Print(result)
}
