package scraper

import (
	"database/sql"

	"scrappy/internal/database"
	"scrappy/internal/types"
)

type Scraper struct {
	scrapees map[string]types.Scrapees
	models   database.Models
}

func NewScraper(db *sql.DB, scrappees map[string]types.Scrapees) *Scraper {
	return &Scraper{
		scrapees: scrappees,
		models:   database.NewModels(db),
	}
}

func (s *Scraper) Run() ([]types.ProductChanged, error) {
	products, err := scrapSites(s.scrapees)
	if err != nil {
		return nil, err
	}
	return detectPriceChanges(products, s.models)
}
