package storage

import (
	"fmt"
	"github.com/pkg/errors"
)

var NameTooShortError = errors.New("name is too short")

const minNameLength = 4

type Product struct {
	code         string
	name         string
	priceHistory PriceHistory
}

func NewProduct(code, name string) (*Product, error) {
	p := Product{
		code:         code,
		priceHistory: *NewPriceHistory(),
	}
	err := p.SetName(name)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Product) SetName(name string) error {
	if len(name) < minNameLength {
		return errors.Wrap(NameTooShortError, name)
	}
	p.name = name
	return nil
}

func (p *Product) AddPriceTimeStamp(stamp *PriceTimeStamp) {
	p.priceHistory.AddPriceTimeStamp(stamp)
}

func (p *Product) String() string {
	return fmt.Sprintf("code: %v, name: %v, %v", p.code, p.name, p.priceHistory.String())
}

func (p *Product) FullHistoryString() string {
	return fmt.Sprintf("code: %v, name: %v\n%v", p.code, p.name, p.priceHistory.FullHistoryString())
}

func (p Product) GetName() string {
	return p.name
}

func (p Product) GetCode() string {
	return p.code
}
