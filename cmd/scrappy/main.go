package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"scrappy/internal/database"
	"scrappy/internal/flags"
	"scrappy/internal/mail"
	"scrappy/internal/scraper"
)

func main() {
	dbPath, configPath, err := flags.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db, err := database.OpenDB(*dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scrapees, err := scraper.GetScrapees(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	changedProducts, err := scraper.NewScraper(db, scrapees).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(changedProducts) < 1 {
		fmt.Println("No price changes detected")
		os.Exit(2)
	}

	m, err := mail.NewMail()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = m.SetSubject("Price Change Detected").SetProducts(changedProducts).Send()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Email sent!")
	os.Exit(0)
}
