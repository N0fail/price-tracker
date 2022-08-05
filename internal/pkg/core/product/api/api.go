package api

import (
	"context"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	ProductList(ctx context.Context) []models.ProductSnapshot
	ProductCreate(ctx context.Context, product models.Product) error
	ProductDelete(ctx context.Context, code string) error
	AddPriceTimeStamp(ctx context.Context, code string, priceTimeStamp models.PriceTimeStamp) error
	FullHistory(ctx context.Context, code string) (models.PriceHistory, error)
}
