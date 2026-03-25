package scraper

import (
	"scrappy/internal/database"
	"scrappy/internal/types"
)

func SaveProducts(products []types.Product, models database.Models) ([]types.Product, error) {
	var cheaperProducts []types.Product

	for _, product := range products {
		siteID, err := models.Site.GetOrCreate(product.Site)
		if err != nil {
			return nil, err
		}

		productID, err := models.Product.GetOrCreate(product.Name, siteID, product.URL)
		if err != nil {
			return nil, err
		}

		_, err = models.History.GetOrCreate(productID, product.Price)
		if err != nil {
			return nil, err
		}

		todays, err := models.History.GetLatest(productID)
		if err != nil {
			return nil, err
		}
		yesterdays, err := models.History.GetPrevious(productID)
		if err != nil {
			return nil, err
		}

		if todays == nil {
			continue
		}

		if yesterdays == nil {
			continue
		}

		if yesterdays.Price > todays.Price {
			cheaperProducts = append(cheaperProducts, product)
		}
	}

	return cheaperProducts, nil
}
