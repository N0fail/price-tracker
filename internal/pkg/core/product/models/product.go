package models

import "fmt"

type Product struct {
	Code string
	Name string
}

func (p Product) String() string {
	return fmt.Sprintf("code: %v, name: %v", p.Code, p.Name)
}
