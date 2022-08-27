package delete

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
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

	err := c.product.ProductDelete(ctx, cmdArgs)
	if err != nil {
		logrus.Error(err.Error())
		return error_codes.GetInternal(err).Error()
	}
	return fmt.Sprintf("product %v was successfully removed", cmdArgs)
}

func (c *command) Name() string {
	return "delete"
}

func (c *command) Help() string {
	return "deletes product from track, args:<code>"
}
