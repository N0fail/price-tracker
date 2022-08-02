package product

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	cachePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache"
	cacheLocalPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache/local"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	Create(product models.Product) error
	Delete(code string) error
	List() []models.ProductSnapshot
	AddPriceTimeStamp(code string, priceTimeStamp models.PriceTimeStamp) error
	FullHistory(code string) (models.PriceHistory, error)
}

func New() Interface {
	return &core{
		cache: cacheLocalPkg.New(),
	}
}

type core struct {
	cache cachePkg.Interface
}

func (c *core) Create(product models.Product) error {
	if len(product.Name) < config.MinNameLength {
		return errors.Wrap(error_codes.ErrNameTooShortError, product.Name)
	}

	return c.cache.ProductCreate(product)
}

func (c *core) Delete(code string) error {
	return c.cache.ProductDelete(code)
}

func (c *core) List() []models.ProductSnapshot {
	return c.cache.ProductList()
}

func (c *core) AddPriceTimeStamp(code string, priceTimeStamp models.PriceTimeStamp) error {
	if priceTimeStamp.Price < 0 {
		return error_codes.ErrNegativePrice
	}
	return c.cache.AddPriceTimeStamp(code, priceTimeStamp)
}

func (c *core) FullHistory(code string) (models.PriceHistory, error) {
	return c.cache.FullHistory(code)
}
