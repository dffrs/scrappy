package main

import (
	"database/sql"
	"fmt"

	"scrappy/internal"
	"scrappy/types"
)

var dbPath = "db.db"

func main() {
	scrapees := map[string]types.Scrapees{
		"simlab":          internal.Simlab{},
		"nextlevelracing": internal.NextLevelRacing{},
		"gtomega":         internal.GTOmega{},
	}

	for fileName, scrapee := range scrapees {
		products, err := scrapee.Run()
		if err != nil {
			panic(err)
		}

		err = internal.SaveAsJSON(products, fileName)
		if err != nil {
			panic(err)
		}
	}

	// NOTE: wip
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mail, err := internal.NewMail()
	if err != nil {
		panic(err)
	}

	subject := "Price Drop Detected"
	message := "Hello there"

	err = mail.SetSubject(&subject).SetMessage(&message).Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent!")
}
