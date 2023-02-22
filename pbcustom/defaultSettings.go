package pbcustom

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/sirupsen/logrus"
)

// AddNewSettingsHandler add a handler to create new default userSettings whenever a user registers
func AddNewSettingsHandler(app *pocketbase.PocketBase, log *logrus.Logger) {
	app.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		log.Println("OnMailerBeforeRecordVerificationSend triggered")
		var countEntries int
		err := app.Dao().DB().
			Select("count(*)").
			From("user_settings").
			Where(dbx.Like("user_id", e.Record.Id)).
			Row(&countEntries)
		if err != nil {
			return err
		}

		log.Printf("found %v entries in database", countEntries)
		if countEntries == 0 {
			log.Println("no entries found")
			collection, err := app.Dao().FindCollectionByNameOrId("user_settings")
			if err != nil {
				return err
			}

			// create new default settings
			record := models.NewRecord(collection)
			record.Set("user_id", e.Record.Id)
			record.Set("clearDoneEntries", false)
			record.Set("bookmarkOrDue", false)
			log.Println("record created")

			// save to database
			err = app.Dao().SaveRecord(record)
			if err != nil {
				return err
			}
			log.Println("record saved")
		}

		return nil
	})
}
