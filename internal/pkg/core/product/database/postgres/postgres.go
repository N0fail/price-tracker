package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	databasePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/database"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"log"
)

func New(pool *pgxpool.Pool) databasePkg.Interface {
	return &postgres{
		pool: pool,
	}
}

type postgres struct {
	pool *pgxpool.Pool
}

func (p *postgres) ProductList(ctx context.Context, page uint32) []models.ProductSnapshot {
	const query = `
	SELECT products.code, products.name, last_price.price, last_price.date
	FROM (SELECT *
		FROM products
		ORDER BY code
		LIMIT $1
		OFFSET $2) as products
	LEFT JOIN
	(SELECT code, price, date FROM (SELECT code,
		price,
		date,
		row_number() over (partition by code order by date desc ) as rank
	FROM price_history) ranks
	WHERE rank < 2) last_price on products.code = last_price.code
	`

	rows, err := p.pool.Query(ctx, query, config.PageSize, page*config.PageSize)
	if err != nil {
		log.Printf("postgres.ProductList: query: %v", err)
		return nil
	}
	defer rows.Close()

	result := make([]models.ProductSnapshot, 0)
	for rows.Next() {
		values, _ := rows.Values()
		var newSnapShot models.ProductSnapshot
		if values[2] != nil {
			err = rows.Scan(&newSnapShot.Code, &newSnapShot.Name, &newSnapShot.LastPrice.Price, &newSnapShot.LastPrice.Date)
		} else { // product with no price
			newSnapShot.Code = values[0].(string)
			newSnapShot.Name = values[1].(string)
			newSnapShot.LastPrice = models.EmptyPriceTimeStamp
		}

		if err != nil {
			log.Printf("postgres.ProductList: scan: %v", err)
			return nil
		}
		result = append(result, newSnapShot)
	}

	return result
}

func (p *postgres) ProductGet(ctx context.Context, code string) (models.Product, error) {
	query, args, err := squirrel.Select("code, name").
		From("products").
		Where(squirrel.Eq{
			"code": code,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return models.Product{}, errors.Errorf("postgres.ProductGet: to sql: %v", err)
	}

	var products []models.Product
	err = pgxscan.Select(ctx, p.pool, &products, query, args...)

	if err != nil {
		return models.Product{}, errors.Errorf("postgres.ProductGet: query: %v", err)
	}

	if len(products) == 0 {
		return models.Product{}, nil
	}

	if len(products) > 1 {
		return models.Product{}, errors.Errorf("postgres.ProductGet: found mutiple producs with code %v", code)
	}

	return products[0], nil
}

func (p *postgres) ProductCreate(ctx context.Context, product models.Product) error {
	existingProduct, err := p.ProductGet(ctx, product.Code)
	if err != nil {
		return err
	}
	if !existingProduct.IsEmpty() {
		return errors.Wrap(error_codes.ErrProductExists, "postgres.ProductCreate")
	}

	query, args, err := squirrel.Insert("products").
		Columns("code, name").
		Values(product.Code, product.Name).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return errors.Errorf("postgres.ProductCreate: to sql: %v", err)
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Errorf("postgres.ProductCreate: to query: %v", err)
	}

	return nil
}

func (p *postgres) ProductDelete(ctx context.Context, code string) error {
	product, err := p.ProductGet(ctx, code)
	if err != nil {
		return err
	}
	if product.IsEmpty() {
		return errors.Wrap(error_codes.ErrProductNotExist, "postgres.ProductDelete")
	}

	query, args, err := squirrel.Delete("products").
		Where(squirrel.Eq{
			"code": code,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return errors.Errorf("postgres.ProductDelete: to sql: %v", err)
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Errorf("postgres.ProductDelete: to query: %v", err)
	}

	return nil
}

func (p *postgres) AddPriceTimeStamp(ctx context.Context, code string, priceTimeStamp models.PriceTimeStamp) error {
	existingProduct, err := p.ProductGet(ctx, code)
	if err != nil {
		return err
	}
	if existingProduct.IsEmpty() {
		return errors.Wrap(error_codes.ErrProductNotExist, "postgres.AddPriceTimeStamp")
	}

	query, args, err := squirrel.Insert("price_history").
		Columns("code, price, date").
		Values(code, priceTimeStamp.Price, priceTimeStamp.Date).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return errors.Errorf("postgres.AddPriceTimeStamp: to sql: %v", err)
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Errorf("postgres.AddPriceTimeStamp: to query: %v", err)
	}

	return nil
}

func (p *postgres) FullHistory(ctx context.Context, code string) (models.PriceHistory, error) {
	existingProduct, err := p.ProductGet(ctx, code)
	if err != nil {
		return nil, err
	}
	if existingProduct.IsEmpty() {
		return nil, errors.Wrap(error_codes.ErrProductNotExist, "postgres.FullHistory")
	}

	query, args, err := squirrel.Select("price, date").
		From("price_history").
		Where(squirrel.Eq{
			"code": code,
		}).
		OrderBy("date").
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, errors.Errorf("postgres.FullHistory: to sql: %v", err)
	}

	var priceHistory models.PriceHistory
	err = pgxscan.Select(ctx, p.pool, &priceHistory, query, args...)
	if err != nil {
		return nil, errors.Errorf("postgres.FullHistory: to query: %v", err)
	}

	return priceHistory, nil
}
