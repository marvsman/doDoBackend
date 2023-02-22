package main

import (
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const TELEGRAM_TOKEN = "6043369518:AAHGeyqy3BmU8buWR38ACj6YzkqwHMOfyXg"
const TELEGRAM_CHAT_ID int64 = -821408765

func main() {
	setupLogging()
	bot, err := initTelegram()
	if err != nil {
		log.Fatal(err)
	}
	app := initPocketbase()

	// check after a record updated
	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		msg := tg.NewMessage(TELEGRAM_CHAT_ID, "OnRecordAfterUpdateRequest")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error(err)
		}

		log.Println(e.Record)

		return nil
	})

	// check before the mail was sent
	app.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		msg := tg.NewMessage(TELEGRAM_CHAT_ID, "BEFORE SEND")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error(err)
		}
		log.Println("BEFORE SEND")
		log.Println(e)

		return nil
	})

	// check after the mail was sent
	app.OnMailerAfterRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		msg := tg.NewMessage(TELEGRAM_CHAT_ID, "AFTER SEND")
		_, err = bot.Send(msg)
		if err != nil {
			log.Error(err)
		}

		log.Println("AFTER SEND")
		log.Println(e)

		return nil
	})

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
	log.SetLevel(logrus.DebugLevel)
}

func initTelegram() (*tg.BotAPI, error) {
	bot, err := tg.NewBotAPI(TELEGRAM_TOKEN)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return bot, nil
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
