package main

import (
	"scrappy/internal"
	"scrappy/types"
)

func main() {
	scrapees := map[string]types.Scrapees{
		"simlab":  internal.Simlab{},
		"gtomega": internal.GTOmega{},
	}

	for fileName, scrapee := range scrapees {
		if err := internal.SaveAsJSON(scrapee.Run(), fileName); err != nil {
			panic(err)
		}
	}
}
