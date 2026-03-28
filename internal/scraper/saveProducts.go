package scraper

import (
	"scrappy/internal/database"
	"scrappy/internal/types"
)

func detectPriceChanges(products []types.Product, models database.Models) ([]types.ProductChanged, error) {
	var cheaperProducts []types.ProductChanged

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

		difference := yesterdays.Price - todays.Price

		if difference < 0 {
			cheaperProducts = append(cheaperProducts, types.ProductChanged{Product: product, Cheaper: false, OldPrice: yesterdays.Price})
		} else if difference > 0 {
			cheaperProducts = append(cheaperProducts, types.ProductChanged{Product: product, Cheaper: true, OldPrice: yesterdays.Price})
		}

	}

	return cheaperProducts, nil
}
