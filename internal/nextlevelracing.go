// Package internal
package internal

import (
	"fmt"
	"strings"

	"scrappy/types"

	"github.com/gocolly/colly/v2"
)

type NextLevelRacing struct {
	site string
	url  string
}

func NewNextLevelRacing() NextLevelRacing {
	return NextLevelRacing{
		site: "https://nextlevelracing.com",
		url:  "https://nextlevelracing.com/racing-cockpits/",
	}
}

func (nlr NextLevelRacing) Run() ([]types.Product, error) {
	nextLevelRacing := NewNextLevelRacing()
	cockspits := make([]types.Product, 0, 20)

	c := colly.NewCollector(
		colly.CacheDir("./cache/nextLevelRacing"),
	)

	c.OnHTML("div[data-category='Racing Cockpits']", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("div[class='s2'] a h3"))
		desc := strings.TrimSpace(e.ChildText("div[class='s2'] div[class='product-edition']"))

		if desc != "" {
			name = fmt.Sprintf("%s - %s", name, desc)
		}

		price := strings.TrimSpace(e.ChildText("div[class='price'] div[class$='eur'] span[class='sale_price'] span:last-child"))
		currency := strings.TrimSpace(e.ChildText("div[class='price'] div[class$='eur'] span[class='sale_price'] span:first-child"))

		if currency != "" {
			price = fmt.Sprintf("%s%s", currency, price)
		}

		cockspits = append(cockspits, types.Product{
			Name:  name,
			Price: price,
			Site:  nextLevelRacing.site,
			URL:   e.ChildAttr("div[class='s2'] a", "href"),
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err := c.Visit(nextLevelRacing.url); err != nil {
		return nil, err
	}

	return cockspits, nil
}
