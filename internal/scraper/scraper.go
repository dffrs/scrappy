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

func NewScraper(db *sql.DB) *Scraper {
	return &Scraper{
		scrapees: map[string]types.Scrapees{
			"simlab":          NewSimlab(),
			"nextlevelracing": NewNextLevelRacing(),
			"gtomega":         NewGTOmega(),
		},
		models: database.NewModels(db),
	}
}

func (s *Scraper) Run() ([]types.Product, error) {
	products, err := scrapSites(s.scrapees)
	if err != nil {
		return nil, err
	}
	return getCheaperProducts(products, s.models)
}
