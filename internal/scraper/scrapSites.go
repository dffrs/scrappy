package scraper

import "scrappy/internal/types"

func ScrapSites(scrapees map[string]types.Scrapees) ([]types.Product, error) {
	p := make([]types.Product, 0)

	for _, scrapee := range scrapees {
		products, err := scrapee.Run()
		if err != nil {
			return nil, err
		}

		p = append(p, products...)
	}

	return p, nil
}
