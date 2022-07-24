package list

import (
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
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
	data := c.product.List()
	if len(data) == 0 {
		return "no products are being tracked now"
	}
	res := make([]string, 0, len(data))
	for _, p := range data {
		res = append(res, p.String())
	}
	return strings.Join(res, "\n")
}

func (c *command) Name() string {
	return "list"
}

func (c *command) Help() string {
	return "get list of all products"
}
