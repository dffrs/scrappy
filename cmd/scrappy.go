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
		products, err := scrapee.Run()
		if err != nil {
			panic(err)
		}

		err = internal.SaveAsJSON(products, fileName)
		if err != nil {
			panic(err)
		}
	}
}
