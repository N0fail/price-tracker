package help

import (
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	"strings"
)

func New(c []commandPkg.Interface) commandPkg.Interface {
	return &command{
		commands: c,
	}
}

type command struct {
	commands []commandPkg.Interface
}

func (c *command) Process(string) string {
	listStr := make([]string, 0, len(c.commands)+2)
	listStr = append(listStr, "you can pass arguments using `"+config.CommandDelimeter+"` delimeter\n example will pass arg1 and arg2 to cmd: /cmd arg1"+config.CommandDelimeter+"arg2")
	for _, command := range c.commands {
		listStr = append(listStr, fmt.Sprintf("/%v: %v", command.Name(), command.Help()))
	}
	listStr = append(listStr, fmt.Sprintf("/%v: %v", c.Name(), c.Help()))
	return strings.Join(listStr, "\n")
}

func (c *command) Name() string {
	return "help"
}

func (c *command) Help() string {
	return "print available commands"
}
