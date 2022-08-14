package list

import (
	"bytes"
	"context"
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"log"
	"strconv"
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

	page, err := strconv.ParseUint(cmdArgs, 10, 32)
	if err != nil {
		log.Print(err.Error())
		return "Error in page format, correct example: 2"
	}
	if page == 0 {
		return "Page number must be positive"
	}
	page = page - 1

	data, err := c.product.ProductList(ctx, uint32(page), config.DefaultResultsPerPage, config.DefaultOrderBy)
	if err != nil {
		log.Print(err.Error())
		return error_codes.ErrExternalProblem.Error()
	}
	var buffer bytes.Buffer
	for _, p := range data {
		buffer.WriteString(p.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (c *command) Name() string {
	return "list"
}

func (c *command) Help() string {
	return fmt.Sprintf("get list of products on given page, there are %v products on one page, args:<page>", config.DefaultResultsPerPage)
}
