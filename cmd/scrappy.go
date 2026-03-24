package main

import (
	"database/sql"
	"fmt"

	"scrappy/internal"
	"scrappy/internal/database"
	"scrappy/types"
)

var dbPath = "./data.db"

func scrapSites(scrapees map[string]types.Scrapees) ([]types.Product, error) {
	p := make([]types.Product, 0, 100)

	for _, scrapee := range scrapees {
		products, err := scrapee.Run()
		if err != nil {
			return nil, err
		}

		p = append(p, products...)
	}

	return p, nil
}

func saveProducts(products []types.Product, models database.Models) ([]types.Product, error) {
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

		// get today's history
		// get yesterday's history
		// if yesterday > today => send email
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

func main() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	scrapees := map[string]types.Scrapees{
		"simlab":          internal.Simlab{},
		"nextlevelracing": internal.NextLevelRacing{},
		"gtomega":         internal.GTOmega{},
	}

	products, err := scrapSites(scrapees)
	if err != nil {
		panic(err)
	}

	models := database.NewModels(db)

	cheaperProducts, err := saveProducts(products, models)
	if err != nil {
		panic(err)
	}

	mail, err := internal.NewMail()
	if err != nil {
		panic(err)
	}

	subject := "Price Drop Detected"

	err = mail.SetSubject(&subject).SetProducts(cheaperProducts).Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent!")
}
