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
