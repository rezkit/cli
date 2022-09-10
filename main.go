package main

import (
	"log"
	"os"

	"github.com/rezkit/cli/app"
)

func main() {

	if err := app.GetApp().Run(os.Args); err != nil {
		log.Fatalln("Error: ", err)
	}
}
