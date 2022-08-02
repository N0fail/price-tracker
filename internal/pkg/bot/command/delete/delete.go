package delete

import (
	"fmt"
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
	err := c.product.Delete(cmdArgs)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("product %v was successfully removed", cmdArgs)
}

func (c *command) Name() string {
	return "delete"
}

func (c *command) Help() string {
	return "deletes product from track, args:<code>"
}
