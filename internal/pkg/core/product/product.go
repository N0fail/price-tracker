package product

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	productApiPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/api"
	cacheLocalPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/cache/local"
	postgresPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/database/postgres"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
)

type Interface interface {
	productApiPkg.Interface
}

func New(pool *pgxpool.Pool) Interface {
	var coreObj core
	// если передали pool используем БД, иначе кэш
	if pool == nil {
		coreObj = core{
			storage: cacheLocalPkg.New(),
		}
	} else {
		coreObj = core{
			storage: postgresPkg.New(pool),
		}
	}
	return &coreObj
}

type core struct {
	storage Interface
}

func (c *core) ProductCreate(ctx context.Context, product models.Product) error {
	return c.storage.ProductCreate(ctx, product)
}

func (c *core) ProductDelete(ctx context.Context, code string) error {
	return c.storage.ProductDelete(ctx, code)
}

func (c *core) ProductList(ctx context.Context, pageNumber, resultsPerPage uint32, orderBy string) ([]models.ProductSnapshot, error) {
	return c.storage.ProductList(ctx, pageNumber, resultsPerPage, orderBy)
}

func (c *core) AddPriceTimeStamp(ctx context.Context, code string, priceTimeStamp models.PriceTimeStamp) error {
	return c.storage.AddPriceTimeStamp(ctx, code, priceTimeStamp)
}

func (c *core) FullHistory(ctx context.Context, code string) (models.PriceHistory, error) {
	return c.storage.FullHistory(ctx, code)
}
