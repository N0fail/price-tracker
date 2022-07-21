package storage

import (
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/config"
	"time"
)

type PriceTimeStamp struct {
	price float64
	date  time.Time
}

func NewPriceTimeStamp(price float64, date time.Time) *PriceTimeStamp {
	return &PriceTimeStamp{
		price: price,
		date:  date,
	}
}

func (p *PriceTimeStamp) SetPrice(newPrice float64) {
	p.price = newPrice
}

func (p *PriceTimeStamp) SetDate(newDate time.Time) {
	p.date = newDate
}

func (p PriceTimeStamp) GetPrice() float64 {
	return p.price
}

func (p PriceTimeStamp) GetDate() time.Time {
	return p.date
}

func (p PriceTimeStamp) String() string {
	return fmt.Sprintf("%v: %v", p.date.Format(config.DateFormat), p.price)
}
