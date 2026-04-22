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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	result := make(map[string]types.Scrapees)

	// Stream the file's content, and decode into a 'Config' object, one at a time
	decoder := json.NewDecoder(f)

	t, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	if t != json.Delim('[') {
		return nil, fmt.Errorf("error reading config: expected array start")
	}

	for decoder.More() {
		var conf Config
		if err := decoder.Decode(&conf); err != nil {
			return nil, err
		}

		result[conf.Name] = conf
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

func (cn Config) getName(e *colly.HTMLElement) string {
	if cn.ProductNamePath == "" {
		return "[Unknown-product-name]"
	}

	name := strings.Split(strings.TrimSpace(e.ChildText(cn.ProductNamePath)), "\n")[0]
	desc := strings.Split(strings.TrimSpace(e.ChildText(cn.ProductDescPath)), "\n")[0]

	// NOTE: Will this be the case for every product ??
	hasDesc := desc != ""
	if hasDesc {
		name = fmt.Sprintf("%s - %s", name, desc)
	}

	return name
}

func (cn Config) getPrice(e *colly.HTMLElement) (float32, error) {
	return extractPrice(strings.TrimSpace(e.ChildText(cn.ProductPricePath)))
}

func (cn Config) getURL(e *colly.HTMLElement) string {
	switch {
	case cn.ProductURLPath == "":
		return cn.Page

	case cn.URLWithSite:
		return fmt.Sprintf("%s%s", cn.Site, e.ChildAttr(cn.ProductURLPath, "href"))

	default:
		return e.ChildAttr(cn.ProductURLPath, "href")
	}
}

func (cn Config) Run() ([]types.Product, error) {
	cockspits := make([]types.Product, 0)

	c := colly.NewCollector()

	if cn.WaitFor != "" {
		htmlContent, err := cn.headlessBrowser()
		if err != nil {
			return nil, err
		}

		c.OnResponse(func(r *colly.Response) {
			r.Body = []byte(*htmlContent)
		})
	}

	c.OnHTML(cn.ContainerPath, func(e *colly.HTMLElement) {
		price, err := cn.getPrice(e)
		if err != nil {
			panic(err)
		}

		cockspits = append(cockspits, types.Product{
			Price: price,
			Site:  cn.Site,
			Name:  cn.getName(e),
			URL:   cn.getURL(e),
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
