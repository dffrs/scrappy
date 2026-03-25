package types

type Product struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Site  string  `json:"site"`
	URL   string  `json:"url"`
}

type Scrapees interface {
	Run() ([]Product, error)
}
