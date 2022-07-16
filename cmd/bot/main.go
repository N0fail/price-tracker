package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.ozon.dev/N0fail/price-tracker/config"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/commander"
	"log"
)

func main() {
	fmt.Println()
	bot, error := tgbotapi.NewBotAPI(config.ApiKey)
	if error != nil {
		log.Panic(error.Error())
	}
	myCommander := commander.Init(bot)
	myCommander.Run()
}
