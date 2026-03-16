package main

import (
	"fmt"

	"scrappy/internal"
	"scrappy/types"
)

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

	// TODO: test to send email. This needs to be removed
	config, err := internal.LoadEnv()
	if err != nil {
		panic(err)
	}

	err = internal.NewMail("Subject: Test\nHello there", config).Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent!")
}
