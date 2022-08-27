package add

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(params[0]) == 0 {
		return error_codes.ErrEmptyCode.Error()
	}

	if len(params[1]) < config.MinNameLength {
		return error_codes.ErrNameTooShortError.Error()
	}

	err := c.product.ProductCreate(ctx, models.Product{
		Code: params[0],
		Name: params[1],
	})
	if err != nil {
		logrus.Error(err.Error())
		return error_codes.ErrExternalProblem.Error()
	}
	return fmt.Sprintf("product %v was successfully added", params[1])
}

func (c *command) Name() string {
	return "add"
}

func (c *command) Help() string {
	return "adds product to track, args:<code>" + config.CommandDelimeter + "<name>"
}
