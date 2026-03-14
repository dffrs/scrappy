package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly/v2"
)

const url = "https://eu.store.bambulab.com/products/a1?from=navigation&id=599117150694776840"

var re = regexp.MustCompile(`(\d)+\sEUR`)

func main() {
	c := colly.NewCollector(
		colly.CacheDir("./cache/bbl"),
	)

	c.OnHTML("span", func(h *colly.HTMLElement) {
		if re.MatchString(h.Text) {
			fmt.Printf("Price: %v\n", h.Text)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
		r.Headers.Set("X-Bbl-Store-Region", "EU")
		r.Headers.Set("X-Bbl-Time-Zone", "Europe/Lisbon")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err := c.Visit(url); err != nil {
		panic(err)
	}
}
