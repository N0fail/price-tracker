package price_history

import (
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
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
	product, err := c.product.Get(cmdArgs)
	if err != nil {
		return err.Error()
	}
	return product.FullHistoryString()
}

func (c *command) Name() string {
	return "price_history"
}

func (c *command) Help() string {
	return "returns price history of a product, args:<code>"
}
