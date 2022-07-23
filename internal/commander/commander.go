package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.ozon.dev/N0fail/price-tracker/config"
	"log"
)

type Commander struct {
	bot *tgbotapi.BotAPI
}

func Init(bot *tgbotapi.BotAPI) *Commander {
	bot.Debug = true
	commander := Commander{bot}
	initHandlers()
	return &commander
}

func (c *Commander) Run() {
	bot := c.bot
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = config.GetUpdatesTimeout

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		var msg tgbotapi.MessageConfig
		if cmd := update.Message.Command(); cmd != "" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, ApplyHandler(cmd, update.Message.CommandArguments()))
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
		}

		bot.Send(msg)
	}
}
