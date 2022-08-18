package database

import (
	"context"
	productApiPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/api"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	productApiPkg.Interface
	ProductGet(ctx context.Context, code string) (models.Product, error)
}
