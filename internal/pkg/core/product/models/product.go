package models

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	Code string `db,json:"code"`
	Name string `db,json:"name"`
}

func (p Product) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Product) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p Product) String() string {
	return fmt.Sprintf("code: %v, name: %v", p.Code, p.Name)
}

func (p Product) IsEmpty() bool {
	return p.Code == ""
}
