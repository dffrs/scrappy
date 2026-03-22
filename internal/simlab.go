// Package internal
package internal

import (
	"fmt"
	"strings"

	"scrappy/types"

	"github.com/gocolly/colly/v2"
)

type Simlab struct {
	site string
	url  string
}

func NewSimlab() Simlab {
	return Simlab{
		site: "https://sim-lab.eu",
		url:  "https://sim-lab.eu/en-pt/collections/sim-racing-cockpits",
	}
}

func (sm Simlab) Run() ([]types.Product, error) {
	simlab := NewSimlab()

	cockspits := make([]types.Product, 0, 20)

	c := colly.NewCollector(
		colly.CacheDir("./cache/simlab"),
	)

	c.OnHTML("ul[id='product-grid'] li[class='grid__item']", func(e *colly.HTMLElement) {
		cockspits = append(cockspits, types.Product{
			Name:  strings.Split(strings.TrimSpace(e.ChildText("h3 a")), "\n")[0],
			Price: strings.TrimSpace(e.ChildText("div[class='price__regular'] span:last-child")),
			Site:  simlab.site,
			URL:   fmt.Sprintf("%s%s", simlab.site, e.ChildAttr("h3 a", "href")),
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err := c.Visit(simlab.url); err != nil {
		return nil, err
	}

	return cockspits, nil
}
