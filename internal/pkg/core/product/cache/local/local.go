package local

import (
	"github.com/pkg/errors"
	cachePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

var (
	ErrProductNotExist = errors.New("product does not exist")
	ErrProductExists   = errors.New("product exists")
)

func New() cachePkg.Interface {
	return &cache{
		data: make(map[string]models.Product, 0),
	}
}

type cache struct {
	data map[string]models.Product
}

func (c *cache) ProductList() []models.Product {
	res := make([]models.Product, 0, len(c.data))
	for _, product := range c.data {
		res = append(res, product)
	}
	return res
}

func (c *cache) ProductGet(code string) (models.Product, error) {
	if _, ok := c.data[code]; !ok {
		return models.Product{}, errors.Wrap(ErrProductNotExist, code)
	}
	return c.data[code], nil
}

func (c *cache) ProductCreate(p models.Product) error {
	if _, ok := c.data[p.Code]; ok {
		return errors.Wrap(ErrProductExists, p.Code)
	}
	c.data[p.Code] = p
	return nil
}

func (c *cache) ProductUpdate(p models.Product) error {
	if _, ok := c.data[p.Code]; !ok {
		return errors.Wrap(ErrProductNotExist, p.String())
	}
	c.data[p.Code] = p
	return nil
}

func (c *cache) ProductDelete(code string) error {
	if _, ok := c.data[code]; !ok {
		return errors.Wrap(ErrProductNotExist, code)
	}
	delete(c.data, code)
	return nil
}

func (c *cache) PriceTimeStampAdd(product models.Product, priceTimeStamp models.PriceTimeStamp) error {
	product.PriceHistory = append(product.PriceHistory, priceTimeStamp)
	return c.ProductUpdate(product)
}
