package product

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	cachePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache"
	cacheLocalPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache/local"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

var (
	ErrNameTooShortError = errors.New("name is too short")
	ErrNegativePrice     = errors.New("price should be positive")
)

type Interface interface {
	Create(product models.Product) error
	Update(product models.Product) error
	Delete(code string) error
	Get(code string) (models.Product, error)
	List() []models.Product
	AddPriceTimeStamp(product models.Product, priceTimeStamp models.PriceTimeStamp) error
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
		return errors.Wrap(ErrNameTooShortError, product.Name)
	}

	return c.cache.ProductCreate(product)
}

func (c *core) Update(product models.Product) error {
	return c.cache.ProductUpdate(product)
}

func (c *core) Delete(code string) error {
	return c.cache.ProductDelete(code)
}

func (c *core) Get(code string) (models.Product, error) {
	return c.cache.ProductGet(code)
}

func (c *core) List() []models.Product {
	return c.cache.ProductList()
}

func (c *core) AddPriceTimeStamp(product models.Product, priceTimeStamp models.PriceTimeStamp) error {
	if priceTimeStamp.Price < 0 {
		return ErrNegativePrice
	}
	return c.cache.PriceTimeStampAdd(product, priceTimeStamp)
}
