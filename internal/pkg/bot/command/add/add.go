package add

import (
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"strings"
)

func New(p productPkg.Interface) commandPkg.Interface {
	return &command{
		product: p,
	}
}

type command struct {
	product productPkg.Interface
}

func (c *command) Process(cmdArgs string) string {
	params := strings.Split(cmdArgs, config.CommandDelimeter)
	if len(params) != 2 {
		return "incorrect number of arguments\n" + c.Help()
	}

	err := c.product.Create(models.Product{
		Code: params[0],
		Name: params[1],
	})
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("product %v was successfully added", params[1])
}

func (c *command) Name() string {
	return "add"
}

func (c *command) Help() string {
	return "adds product to track, args:<code>" + config.CommandDelimeter + "<name>"
}
