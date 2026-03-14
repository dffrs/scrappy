// Package types
package types

type Product interface {
	Name() string
	Price() string
}

type Scrapees interface {
	Run() []Product
}
