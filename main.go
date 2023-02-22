package main

import (
	"fmt"
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
		msg := tg.NewMessage(TELEGRAM_CHAT_ID, fmt.Sprintf("Updated record: %v", e.Record.Id))
		_, err = bot.Send(msg)
		if err != nil {
			log.Error(err)
		}

		log.Println(e.Record)

		return nil
	})

	// check before the mail was sent
	app.OnMailerBeforeRecordVerificationSend().Add(func(e *core.MailerRecordEvent) error {
		_ = sendMsg(bot, "Before send")
		log.Println("BEFORE SEND")

		var countEntries int
		err := app.Dao().DB().
			Select("count(*)").
			From("userSettings").
			Where(dbx.Like("user_id", e.Record.Id)).
			Row(&countEntries)
		if err != nil {
			return err
		}

		_ = sendMsg(bot, fmt.Sprintf("Found %v entries for this user", countEntries))
		log.Printf("Found %v entries for this user", countEntries)

		if countEntries == 0 {
			_ = sendMsg(bot, "no entry found")
			collection, err := app.Dao().FindCollectionByNameOrId("userSettings")
			if err != nil {
				return err
			}

			_ = sendMsg(bot, fmt.Sprintf("Collection found: %s", collection.Name))

			record := models.NewRecord(collection)
			record.Set("user_id", e.Record.Id)
			record.Set("clearDoneEntries", false)
			record.Set("bookmarkOrDue", false)

			_ = sendMsg(bot, fmt.Sprintf("Created new settings: %s", record.Id))
		}

		return nil
	})

	// start pocketbase
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
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

func sendMsg(bot *tg.BotAPI, message string) error {
	msg := tg.NewMessage(TELEGRAM_CHAT_ID, message)
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func initPocketbase() *pocketbase.PocketBase {
	pocketbase.Version = "v0.0.1"
	app := pocketbase.New()
	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: true,
	})

	return app
}
