package storage

import (
	"fmt"
	"time"
)

const DateFormat = "2 Jan 2006 15:04"

type PriceTimeStamp struct {
	price uint64
	date  time.Time
}

func NewPriceTimeStamp(price uint64, date time.Time) *PriceTimeStamp {
	return &PriceTimeStamp{
		price: price,
		date:  date,
	}
}

func (p *PriceTimeStamp) SetPrice(newPrice uint64) {
	p.price = newPrice
}

func (p *PriceTimeStamp) SetDate(newDate time.Time) {
	p.date = newDate
}

func (p PriceTimeStamp) GetPrice() uint64 {
	return p.price
}

func (p PriceTimeStamp) GetDate() time.Time {
	return p.date
}

func (p PriceTimeStamp) String() string {
	return fmt.Sprintf("%v: %v", p.date.Format(DateFormat), p.price)
}
