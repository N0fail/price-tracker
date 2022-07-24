package cache

import (
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	ProductList() []models.Product
	ProductCreate(p models.Product) error
	ProductGet(code string) (models.Product, error)
	ProductDelete(code string) error
	ProductUpdate(p models.Product) error
	PriceTimeStampAdd(product models.Product, priceTimeStamp models.PriceTimeStamp) error
}
