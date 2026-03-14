// Package internal
package internal

type Simlab struct {
	url   string
	Name  string
	Price string
}

func NewSimlab() *Simlab {
	return &Simlab{
		url: "https://sim-lab.eu/en-pt/collections/sim-racing-cockpits",
	}
}

func (sm *Simlab) URL() string {
	return sm.url
}
