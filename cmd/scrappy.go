package main

import (
	"fmt"

	"scrappy/internal"
	"scrappy/types"
)

func main() {
	scrapees := []types.Scrapees{
		internal.Simlab{},
		internal.GTOmega{},
	}

	for _, scrapee := range scrapees {
		for _, product := range scrapee.Run() {
			fmt.Printf("\n%s\n%s\n", product.Name(), product.Price())
		}
	}
}
