package commander

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/N0fail/price-tracker/config"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/storage"
	"strconv"
	"strings"
	"time"
)

type commandSignature func(string) string
type Handler struct {
	name  string
	help  string
	apply commandSignature
}

func (h Handler) String() string {
	return fmt.Sprintf("/%v: %v", h.name, h.help)
}

var HandlerExists = errors.New("handler exists")
var handlers = make(map[string]Handler, 0)

func initHandlers() {
	NewHandler("help", "print available commands", func(string) string {
		listStr := make([]string, 0, len(handlers)+1)
		listStr = append(listStr, "you can pass arguments using `"+config.CommandDelimeter+"` delimeter\n example will pass arg1 and arg2 to cmd: /cmd arg1"+config.CommandDelimeter+"arg2")
		for _, handler := range handlers {
			listStr = append(listStr, handler.String())
		}
		return strings.Join(listStr, "\n")
	})

	NewHandler("list", "get list of all products", func(cmdArgs string) string {
		data := storage.List()
		if len(data) == 0 {
			return "no products are being tracked now"
		}
		res := make([]string, 0, len(data))
		for _, p := range data {
			res = append(res, p.String())
		}
		return strings.Join(res, "\n")
	})

	NewHandler("add", "adds product to track, args:<code>"+config.CommandDelimeter+"<name>", func(cmdArgs string) string {
		params := strings.Split(cmdArgs, config.CommandDelimeter)
		if len(params) != 2 {
			return "incorrect number of arguments"
		}
		p, err := storage.NewProduct(params[0], params[1])
		if err != nil {
			return err.Error()
		}
		err = storage.AddProduct(p)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("product %v was successfully added", p.GetName())
	})

	NewHandler("delete", "deletes product from track, args:<code>", func(cmdArgs string) string {
		err := storage.DeleteProduct(cmdArgs)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("product %v was successfully removed", cmdArgs)
	})

	NewHandler("add_price", "adds price to product track, args:<code>"+config.CommandDelimeter+"<date>"+config.CommandDelimeter+"<price>\ndate format: "+config.DateFormat, func(cmdArgs string) string {
		params := strings.Split(cmdArgs, config.CommandDelimeter)
		if len(params) != 3 {
			return "incorrect number of arguments"
		}
		code, date, price := params[0], params[1], params[2]
		product, err := storage.GetProduct(code)
		if err != nil {
			return err.Error()
		}
		dateTime, err := time.Parse(config.DateFormat, date)
		if err != nil {
			return err.Error()
		}
		priceUint, err := strconv.ParseUint(price, 10, 64)
		if err != nil {
			return err.Error()
		}
		stamp := storage.NewPriceTimeStamp(priceUint, dateTime)
		product.AddPriceTimeStamp(stamp)
		return fmt.Sprintf("Price %v was successfully added for product %v", stamp.String(), product.GetName())
	})

	NewHandler("price_history", "returns price history of a product, args:<code>", func(cmdArgs string) string {
		product, err := storage.GetProduct(cmdArgs)
		if err != nil {
			return err.Error()
		}
		return product.FullHistoryString()
	})
}

func NewHandler(name, help string, apply commandSignature) error {
	if _, ok := handlers[name]; ok {
		return errors.Wrap(HandlerExists, name)
	}
	handlers[name] = Handler{
		name:  name,
		help:  help,
		apply: apply,
	}
	return nil
}

func ApplyHandler(name, cmdArgs string) string {
	if _, ok := handlers[name]; !ok {
		return "unknown command"
	}
	return handlers[name].apply(cmdArgs)
}
