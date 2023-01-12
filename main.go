package main

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase"
)

func main() {
	pocketbase.Version = "v0.0.1"

	app := pocketbase.New()
	if err := app.Start(); err != nil {
		fmt.Println("bin hier")
		log.Fatal(err)
	}
}
