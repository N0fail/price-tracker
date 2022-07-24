package models

import (
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"time"
)

type PriceTimeStamp struct {
	Price float64
	Date  time.Time
}

func (p PriceTimeStamp) String() string {
	return fmt.Sprintf("%v: %v", p.Date.Format(config.DateFormat), p.Price)
}
