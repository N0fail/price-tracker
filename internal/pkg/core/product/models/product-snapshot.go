package models

import (
	"encoding/json"
	"fmt"
)

type ProductSnapshot struct {
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	LastPrice PriceTimeStamp `json:"last_price"`
}

type ProductSnapshots []ProductSnapshot

func (p ProductSnapshot) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProductSnapshot) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p ProductSnapshots) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProductSnapshots) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p ProductSnapshot) String() string {
	return fmt.Sprintf("code: %v, name: %v, last price: %v", p.Code, p.Name, p.LastPrice)
}
