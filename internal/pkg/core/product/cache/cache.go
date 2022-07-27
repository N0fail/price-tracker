package cache

import (
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	ProductList() []models.ProductSnapshot
	ProductCreate(p models.Product) error
	ProductDelete(code string) error
	AddPriceTimeStamp(code string, priceTimeStamp models.PriceTimeStamp) error
	FullHistory(code string) (models.PriceHistory, error)
}
