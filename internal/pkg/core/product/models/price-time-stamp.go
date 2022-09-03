package models

import (
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"time"
)

type PriceTimeStamp struct {
	Price float64   `db,json:"price"`
	Date  time.Time `db,json:"date"`
}

func (p PriceTimeStamp) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PriceTimeStamp) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
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
