package models

import "fmt"

type ProductSnapshot struct {
	Code      string
	Name      string
	LastPrice PriceTimeStamp
}

func (p ProductSnapshot) String() string {
	return fmt.Sprintf("code: %v, name: %v, last price: %v", p.Code, p.Name, p.LastPrice)
}
