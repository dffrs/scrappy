package main

import (
	"fmt"

	"scrappy/internal"
	"scrappy/types"
)

func main() {
	temp := []types.Product[internal.Simlab]{
		&internal.Simlab{},
	}

	for _, p := range temp {
		for _, v := range p.Run() {
			fmt.Printf("\n%s\n%s\n", v.Name, v.Price)
		}
	}
}
