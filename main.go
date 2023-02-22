package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
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

	// create default userSettings entry for new users
	app.OnRecordBeforeConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
		collection, err := app.Dao().FindCollectionByNameOrId("userSettings")
		if err != nil {
			return err
		}

		record := models.NewRecord(collection)
		record.Set("userID", e.Record.Id)
		record.Set("clearDoneEntries", false)
		record.Set("bookmarkOrDue", false)

		return nil
	})
}
