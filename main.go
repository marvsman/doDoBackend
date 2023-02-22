package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	pocketbase.Version = "v0.0.1"

	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	//app.OnRecordBeforeConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
	//	collection, err := app.Dao().FindCollectionByNameOrId("userSettings")
	//	if err != nil {
	//		return err
	//	}
	//
	//	record := models.NewRecord(collection)
	//	record.Set("title", "Lorem ipsum")
	//	record.Set("active", true)
	//	record.Set("someOtherField", 123)
	//	log.Println(e.Record)
	//})
}
