package main

import (
	"encoding/json"
	"log"

	"github.com/dougrich/tracerygo"
)

func main() {
	rawg := make(tracerygo.RawGrammar)
	err := json.Unmarshal([]byte(`{"move":["flock","race","glide","dance","flee","lie"],"bird":["swan","heron","sparrow","swallow","wren","robin"],"agent":["cloud","wave","#bird#","boat","ship"],"transVerb":["forget","plant","greet","remember","embrace","feel","love"],"emotion":["sorrow","gladness","joy","heartache","love","forgiveness","grace"],"substance":["#emotion#","mist","fog","glass","silver","rain","dew","cloud","virtue","sun","shadow","gold","light","darkness"],"adj":["fair","bright","splendid","divine","inseparable","fine","lazy","grand","slow","quick","graceful","grave","clear","faint","dreary"],"doThing":["come","move","cry","weep","laugh","dream"],"verb":["fleck","grace","bless","dapple","touch","caress","smooth","crown","veil"],"ground":["glen","river","vale","sea","meadow","forest","glade","grass","sky","waves"],"poeticAdj":["#substance#-#verb.ed#"],"poeticDesc":["#poeticAdj#","by #substance# #verb#'d","#adj# with #substance#","#verb.ed# with #substance#"],"ah":["ah","alas","oh","yet","but","and"],"on":["on","in","above","beneath","under","by"],"punctutation":[",",":"," ","!",".","?"],"noun":["#ground#","#agent#"],"line":["My #noun#, #poeticDesc#, my #adj# one","More #adj# than #noun# #poeticDesc#","#move.capitalize# with me #on# #poeticAdj# #ground#","The #agent.s# #move#, #adj# and #adj#","#poeticDesc.capitalize#, #poeticDesc#, #ah#, #poeticDesc#","How #adj# is the #poeticDesc# #sub#","#poeticDesc.capitalize# with #emotion#, #transVerb.s# the #noun#"],"poem":["#line##punctutation#<br>#line##punctutation#<br>#line##punctutation#<br>#line#."],"origin":"#[sub:#noun#]poem#"}`), &rawg)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON %v", err)
	}

	result, err := rawg.Evaluate("origin", 0, 0)
	if err != nil {
		log.Fatalf("Error directly evaluating %v", err)
	}

	log.Print(result)
}
