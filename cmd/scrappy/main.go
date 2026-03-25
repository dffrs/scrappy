package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"scrappy/internal/database"
	"scrappy/internal/mail"
	"scrappy/internal/scraper"
	"scrappy/internal/types"
)

func main() {
	dbPath := flag.String("db", "./data.db", "path to sqlite database")
	flag.Parse()

	db, err := database.OpenDB(*dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scrapees := map[string]types.Scrapees{
		"simlab":          scraper.NewSimlab(),
		"nextlevelracing": scraper.NewNextLevelRacing(),
		"gtomega":         scraper.NewGTOmega(),
	}

	products, err := scraper.ScrapSites(scrapees)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	models := database.NewModels(db)

	cheaperProducts, err := scraper.SaveProducts(products, models)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(cheaperProducts) < 1 {
		fmt.Println("No price drop detected")
		os.Exit(2)
	}

	m, err := mail.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = m.SetSubject("Price Drop Detected").SetProducts(cheaperProducts).Send()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Email sent!")
	os.Exit(0)
}
