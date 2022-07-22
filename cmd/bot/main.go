package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/commander"
	"log"
	"os"
)

func main() {
	apiKey, exists := os.LookupEnv("PriceTrackerApiKey")
	if !exists {
		log.Fatal(errors.New("PriceTrackerApiKey environment variable expected\nuse `export PriceTrackerApiKey=` to set your bot ApiKey"))
	}
	bot, error := tgbotapi.NewBotAPI(apiKey)
	if error != nil {
		log.Fatal(error.Error())
	}
	myCommander := commander.Init(bot)
	myCommander.Run()
}
