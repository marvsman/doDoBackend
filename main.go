package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
)

func main() {
	pocketbase.Version = "v0.0.1"

	app := pocketbase.New()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
