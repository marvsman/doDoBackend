package main

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// CreateNewSettingsHandler add a handler to create new default userSettings whenever a user registers
func CreateNewSettingsHandler(app *pocketbase.PocketBase) {
	app.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		var countEntries int
		err := app.Dao().DB().
			Select("count(*)").
			From("user_settings").
			Where(dbx.Like("user_id", e.Record.Id)).
			Row(&countEntries)
		if err != nil {
			return err
		}

		if countEntries == 0 {
			collection, err := app.Dao().FindCollectionByNameOrId("userSettings")
			if err != nil {
				return err
			}

			// create new default settings
			record := models.NewRecord(collection)
			record.Set("user_id", e.Record.Id)
			record.Set("clearDoneEntries", false)
			record.Set("bookmarkOrDue", false)

			// save to database
			err = app.Dao().SaveRecord(record)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
