package models

import "fmt"

type Product struct {
	Code         string
	Name         string
	PriceHistory PriceHistory
}

func (p Product) String() string {
	return fmt.Sprintf("code: %v, name: %v, %v", p.Code, p.Name, p.PriceHistory.String())
}

func (p Product) FullHistoryString() string {
	return fmt.Sprintf("code: %v, name: %v\n%v", p.Code, p.Name, p.PriceHistory.FullHistoryString())
}
