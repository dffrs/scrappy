package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scrappy/types"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

type GTOmega struct {
	site string
	url  string
}

func NewGTOmega() GTOmega {
	return GTOmega{
		site: "https://www.gtomega.eu",
		url:  "https://www.gtomega.eu/collections/cockpits",
	}
}

func headlessBrowser(url string) (*string, error) {
	var htmlContent string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("ul[id='gf-products']"),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, err
	}

	return &htmlContent, nil
}

func (gto GTOmega) Run() ([]types.Product, error) {
	gtomega := NewGTOmega()
	cockspits := make([]types.Product, 0, 20)

	htmlContent, err := headlessBrowser(gtomega.url)
	if err != nil {
		return nil, err
	}

	c := colly.NewCollector(
		colly.CacheDir("./cache/gtomega"),
	)

	c.OnResponse(func(r *colly.Response) {
		r.Body = []byte(*htmlContent)
	})

	c.OnHTML("ul[id='gf-products'] div[class='spf-product__info']", func(e *colly.HTMLElement) {
		price, err := extractPrice(strings.TrimSpace(e.ChildText("div[class='spf-product-card__price-wrapper'] span:last-child")))
		if err != nil {
			panic(err)
		}

		cockspits = append(cockspits, types.Product{
			Name:  strings.TrimSpace(e.ChildText("a")),
			Price: price,
			Site:  gtomega.site,
			URL:   fmt.Sprintf("%s%s", gtomega.site, e.ChildAttr("a", "href")),
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err := c.Visit(gtomega.url); err != nil {
		return nil, err
	}

	return cockspits, nil
}
