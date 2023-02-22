package main

import (
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	setupLogging()
	app := initPocketbase()

	// create default userSettings entry for new users
	app.OnRecordBeforeConfirmVerificationRequest().Add(func(e *core.RecordConfirmVerificationEvent) error {
		var countEntries int
		err := app.Dao().DB().
			Select("count(*)").
			From("userSettings").
			Where(dbx.Like("user_id", e.Record.Id)).
			Row(&countEntries)
		if err != nil {
			return err
		}

		log.Printf("Found %v entries for this user", countEntries)

		if countEntries == 0 {
			collection, err := app.Dao().FindCollectionByNameOrId("userSettings")
			if err != nil {
				return err
			}

			record := models.NewRecord(collection)
			record.Set("user_id", e.Record.Id)
			record.Set("clearDoneEntries", false)
			record.Set("bookmarkOrDue", false)
		}

		return nil
	})

	startPocketbase(app)
}

func setupLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func initPocketbase() *pocketbase.PocketBase {
	pocketbase.Version = "v0.0.1"
	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	return app
}

func startPocketbase(app *pocketbase.PocketBase) {
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
