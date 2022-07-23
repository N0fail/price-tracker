package storage

import "github.com/pkg/errors"

var data map[string]*Product
var ProductNotExist = errors.New("product does not exist")
var ProductExists = errors.New("product exists")

func init() {
	data = make(map[string]*Product, 0)
}

func List() []*Product {
	res := make([]*Product, 0, len(data))
	for _, product := range data {
		res = append(res, product)
	}
	return res
}

func GetProduct(code string) (*Product, error) {
	if _, ok := data[code]; !ok {
		return nil, errors.Wrap(ProductNotExist, code)
	}
	return data[code], nil
}

func AddProduct(p *Product) error {
	if _, ok := data[p.GetCode()]; ok {
		return errors.Wrap(ProductExists, p.GetCode())
	}
	data[p.GetCode()] = p
	return nil
}

func UpdateProduct(p *Product) error {
	if _, ok := data[p.GetCode()]; !ok {
		return errors.Wrap(ProductNotExist, p.String())
	}
	data[p.GetCode()] = p
	return nil
}

func DeleteProduct(code string) error {
	if _, ok := data[code]; !ok {
		return errors.Wrap(ProductNotExist, code)
	}
	delete(data, code)
	return nil
}
