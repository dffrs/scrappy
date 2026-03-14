// Package types
package types

type Product[T any] interface {
	Run() []T
}
