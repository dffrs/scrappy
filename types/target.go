// Package types
package types

type Product struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Scrapees interface {
	Run() []Product
}
