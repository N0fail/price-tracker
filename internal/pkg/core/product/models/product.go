package models

import "fmt"

type Product struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

func (p Product) String() string {
	return fmt.Sprintf("code: %v, name: %v", p.Code, p.Name)
}

func (p Product) IsEmpty() bool {
	return p.Code == ""
}
