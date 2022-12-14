//go:generate mockgen -source ./api.go -destination ./mocks/api.go -package=mock_api
package api

import (
	"context"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	ProductList(ctx context.Context, pageNumber, resultsPerPage uint32, orderBy string) (models.ProductSnapshots, error)
	ProductCreate(ctx context.Context, product models.Product) error
	ProductDelete(ctx context.Context, code string) error
	PriceTimeStampAdd(ctx context.Context, code string, priceTimeStamp models.PriceTimeStamp) error
	PriceHistory(ctx context.Context, code string) (models.PriceHistory, error)
}
