package price_history

import (
	"context"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"log"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	history, err := c.product.FullHistory(ctx, cmdArgs)
	if err != nil {
		log.Println(err.Error())
		return error_codes.GetInternal(err).Error()
	}
	return history.String()
}

func (c *command) Name() string {
	return "price_history"
}

func (c *command) Help() string {
	return "returns price history of a product, args:<code>"
}
