package main

import (
	"fmt"

	"scrappy/internal"
	"scrappy/types"
)

func main() {
	// TODO: Remove me
	if false {
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
	}

	mail, err := internal.NewMail()
	if err != nil {
		panic(err)
	}

	subject := "Test with template"
	message := "Hello there"

	err = mail.SetSubject(&subject).SetMessage(&message).Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent!")
}
