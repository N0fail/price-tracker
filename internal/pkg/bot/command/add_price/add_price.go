package add_price

import (
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	commandPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"log"
	"strconv"
	"strings"
	"time"
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
	if len(params) != 3 {
		return "incorrect number of arguments\n" + c.Help()
	}
	code, date, price := params[0], params[1], params[2]
	dateTime, err := time.Parse(config.DateFormat, date)
	if err != nil {
		log.Print(err.Error())
		return "Error in date format, correct example: " + config.DateFormat
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Print(err.Error())
		return "Error in price format, correct example: 123.45"
	}
	priceTimeStamp := models.PriceTimeStamp{
		Price: priceFloat,
		Date:  dateTime,
	}
	err = c.product.AddPriceTimeStamp(code, priceTimeStamp)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Price %v was successfully added for product %v", priceTimeStamp.String(), code)
}

func (c *command) Name() string {
	return "add_price"
}

func (c *command) Help() string {
	return "adds price to product track, args:<code>" + config.CommandDelimeter + "<date>" + config.CommandDelimeter + "<price>\ndate format: " + config.DateFormat
}
