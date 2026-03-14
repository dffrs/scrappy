// Package internal
package internal

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

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

func (sm *Simlab) Run() []Simlab {
	cockspits := make([]Simlab, 0, 20)

	c := colly.NewCollector(
		colly.CacheDir("./cache/simlab"),
	)

	c.OnHTML("ul[id='product-grid'] li[class='grid__item']", func(e *colly.HTMLElement) {
		cockpit := NewSimlab()

		cockpit.Name = strings.Split(strings.TrimSpace(e.ChildText("h3 a")), "\n")[0]
		cockpit.Price = strings.TrimSpace(e.ChildText("div[class='price__regular'] span:last-child"))

		cockspits = append(cockspits, *cockpit)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err := c.Visit(NewSimlab().url); err != nil {
		panic(err)
	}

	return cockspits
}
