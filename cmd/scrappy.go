package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"scrappy/internal"
	"scrappy/internal/database"
	"scrappy/types"
)

var dbPath = "./data.db"

func main() {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	scrapees := map[string]types.Scrapees{
		"simlab":          internal.Simlab{},
		"nextlevelracing": internal.NextLevelRacing{},
		"gtomega":         internal.GTOmega{},
	}

	products, err := internal.ScrapSites(scrapees)
	if err != nil {
		panic(err)
	}

	models := database.NewModels(db)

	cheaperProducts, err := internal.SaveProducts(products, models)
	if err != nil {
		panic(err)
	}

	if len(cheaperProducts) < 1 {
		fmt.Println("No price drop detected")
		os.Exit(2)
	}

	mail, err := internal.NewMail()
	if err != nil {
		panic(err)
	}

	err = mail.SetSubject("Price Drop Detected").SetProducts(cheaperProducts).Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent!")
	os.Exit(0)
}
