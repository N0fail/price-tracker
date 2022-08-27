package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	"log"
	"os"
)

var ErrHandlerExists = errors.New("handler exists")

type Interface interface {
	RegisterCommand(cmd commandPkg.Interface) error
	GetCommands() []commandPkg.Interface
	Run()
}

func MustNew() Interface {
	apiKey, exists := os.LookupEnv("PriceTrackerApiKey")
	if !exists {
		log.Fatal(errors.New("PriceTrackerApiKey environment variable expected\nuse `export PriceTrackerApiKey=` to set your bot ApiKey"))
	}
	bot, error := tgbotapi.NewBotAPI(apiKey)
	if error != nil {
		log.Fatal(error.Error())
	}
	return &commander{
		bot:      bot,
		commands: make(map[string]commandPkg.Interface, 0),
	}
}

type commander struct {
	bot      *tgbotapi.BotAPI
	commands map[string]commandPkg.Interface
}

func (c *commander) RegisterCommand(cmd commandPkg.Interface) error {
	if _, ok := c.commands[cmd.Name()]; ok {
		return errors.Wrap(ErrHandlerExists, cmd.Name())
	}
	c.commands[cmd.Name()] = cmd
	return nil
}

func (c *commander) GetCommands() []commandPkg.Interface {
	result := make([]commandPkg.Interface, 0, len(c.commands))
	for _, command := range c.commands {
		result = append(result, command)
	}
	return result
}

func (c *commander) Run() {
	bot := c.bot
	logrus.Info("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = config.GetUpdatesTimeout

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		logrus.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)
		var msg tgbotapi.MessageConfig
		if cmd := update.Message.Command(); cmd != "" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, c.applyCommand(cmd, update.Message.CommandArguments()))
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
		}

		bot.Send(msg)
	}
}

func (c *commander) applyCommand(name, cmdArgs string) string {
	if _, ok := c.commands[name]; !ok {
		return "unknown command"
	}
	return c.commands[name].Process(cmdArgs)
}
