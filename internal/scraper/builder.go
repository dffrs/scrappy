package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"scrappy/internal/types"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

type Config struct {
	Name             string `json:"name"`
	Site             string `json:"site"`
	Page             string `json:"page"`
	ContainerPath    string `json:"containerPath"`
	ProductNamePath  string `json:"productNamePath"`
	ProductDescPath  string `json:"productDescPath"`
	ProductPricePath string `json:"productPricePath"`
	ProductURLPath   string `json:"productURLPath"`
	URLWithSite      bool   `json:"urlWithSite"`
	WaitFor          string `json:"waitFor"`
}

func GetScrapees(path string) (map[string]types.Scrapees, error) {
	// open file
	f, err := os.ReadFile(path) // TODO: use scanner
	if err != nil {
		return nil, err
	}

	var content []Config
	err = json.Unmarshal(f, &content)
	if err != nil {
		return nil, err
	}

	result := make(map[string]types.Scrapees)
	for _, scr := range content {
		result[scr.Name] = scr
	}

	return result, nil
}

func (cn Config) headlessBrowser() (*string, error) {
	var htmlContent string

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(cn.Page),
		chromedp.WaitVisible(cn.WaitFor),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, err
	}

	return &htmlContent, nil
}

func (cn Config) Run() ([]types.Product, error) {
	cockspits := make([]types.Product, 0)

	c := colly.NewCollector(
		colly.CacheDir("./cache/nextLevelRacing"),
	)

	if cn.WaitFor != "" {
		htmlContent, err := cn.headlessBrowser()
		if err != nil {
			return nil, err
		}

		c.OnResponse(func(r *colly.Response) {
			r.Body = []byte(*htmlContent)
		})
	}

	// TODO: have fallbacks for every property
	c.OnHTML(cn.ContainerPath, func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText(cn.ProductNamePath))
		desc := strings.TrimSpace(e.ChildText(cn.ProductDescPath))

		if desc != "" {
			name = fmt.Sprintf("%s - %s", name, desc)
		}

		price, err := extractPrice(strings.TrimSpace(e.ChildText(cn.ProductPricePath)))
		if err != nil {
			panic(err)
		}

		var url string
		if cn.ProductURLPath == "" {
			url = cn.Page
		} else if cn.URLWithSite {
			url = fmt.Sprintf("%s%s", cn.Site, e.ChildAttr(cn.ProductURLPath, "href"))
		} else {
			url = e.ChildAttr(cn.ProductURLPath, "href")
		}

		cockspits = append(cockspits, types.Product{
			Name:  name,
			Price: price,
			Site:  cn.Site,
			URL:   url,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())

		r.Headers.Set("User-Agent", "Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/113.0 Firefox/113.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	if err := c.Visit(cn.Page); err != nil {
		return nil, err
	}

	return cockspits, nil
}
