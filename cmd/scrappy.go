package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	i "scrappy/internal"
)

func main() {
	cockspits := make([]*i.Simlab, 0, 20)

	c := colly.NewCollector(
		colly.CacheDir("./cache/simlab"),
	)

	c.OnHTML("ul[id='product-grid'] li[class='grid__item']", func(e *colly.HTMLElement) {
		cockpit := i.NewSimlab()

		cockpit.Name = strings.Split(strings.TrimSpace(e.ChildText("h3 a")), "\n")[0]
		cockpit.Price = strings.TrimSpace(e.ChildText("div[class='price__regular'] span:last-child"))

		cockspits = append(cockspits, cockpit)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		for _, v := range cockspits {
			fmt.Printf("\ncockpit: %s\n%s\n", v.Name, v.Price)
		}
	})

	if err := c.Visit(i.NewSimlab().URL()); err != nil {
		panic(err)
	}
}
