package models

import (
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"time"
)

type PriceTimeStamp struct {
	Price float64   `db:"price"`
	Date  time.Time `db:"date"`
}

func (p PriceTimeStamp) String() string {
	if p.IsEmpty() {
		return "no price"
	}
	return fmt.Sprintf("%v: %v", p.Date.Format(config.DateFormat), p.Price)
}

var EmptyPriceTimeStamp = PriceTimeStamp{
	Price: -1,
}

func (p PriceTimeStamp) IsEmpty() bool {
	return p.Price < 0
}
