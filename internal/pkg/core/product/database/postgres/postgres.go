package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	databasePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/database"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"time"
)

func New(pool DbConn) databasePkg.Interface {
	return &postgres{
		pool: pool,
	}
}

type postgres struct {
	pool DbConn
}

type DbConn interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

func (p *postgres) ProductList(ctx context.Context, pageNumber, resultsPerPage uint32, orderBy string) ([]models.ProductSnapshot, error) {
	const query = `
	SELECT products.code, products.name, last_price.price, last_price.date
	FROM (SELECT *
		FROM products
		ORDER BY $3
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

	rows, err := p.pool.Query(ctx, query, resultsPerPage, pageNumber*resultsPerPage, orderBy)
	if err != nil {
		return nil, errors.Wrapf(err, "postgres.ProductList: query")
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, error_codes.ErrNoEntries
	}

	result := make([]models.ProductSnapshot, 0)
	for {
		values, _ := rows.Values()
		var newSnapShot models.ProductSnapshot
		if values[2] != nil {
			newSnapShot.Code = values[0].(string)
			newSnapShot.Name = values[1].(string)
			newSnapShot.LastPrice = models.PriceTimeStamp{
				Price: values[2].(float64),
				Date:  values[3].(time.Time),
			}
		} else { // product with no price
			newSnapShot.Code = values[0].(string)
			newSnapShot.Name = values[1].(string)
			newSnapShot.LastPrice = models.EmptyPriceTimeStamp
		}

		result = append(result, newSnapShot)

		if !rows.Next() {
			break
		}
	}

	return result, nil
}

func (p *postgres) ProductGet(ctx context.Context, code string) (models.Product, error) {
	query, args, err := squirrel.Select("code, name").
		From("products").
		Where(squirrel.Eq{
			"code": code,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return models.Product{}, errors.Wrapf(err, "postgres.ProductGet: to sql")
	}

	var products []models.Product
	err = pgxscan.Select(ctx, p.pool, &products, query, args...)

	if err != nil {
		return models.Product{}, errors.Wrapf(err, "postgres.ProductGet: query")
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
		return errors.Wrapf(error_codes.ErrProductExists, "postgres.ProductCreate")
	}

	query, args, err := squirrel.Insert("products").
		Columns("code, name").
		Values(product.Code, product.Name).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return errors.Wrapf(err, "postgres.ProductCreate: to sql")
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "postgres.ProductCreate: to query:")
	}

	return nil
}

func (p *postgres) ProductDelete(ctx context.Context, code string) error {
	product, err := p.ProductGet(ctx, code)
	if err != nil {
		return err
	}
	if product.IsEmpty() {
		return errors.Wrapf(error_codes.ErrProductNotExist, "postgres.ProductDelete")
	}

	query, args, err := squirrel.Delete("products").
		Where(squirrel.Eq{
			"code": code,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return errors.Wrapf(err, "postgres.ProductDelete: to sql")
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "postgres.ProductDelete: to query")
	}

	return nil
}

func (p *postgres) PriceTimeStampAdd(ctx context.Context, code string, priceTimeStamp models.PriceTimeStamp) error {
	existingProduct, err := p.ProductGet(ctx, code)
	if err != nil {
		return err
	}
	if existingProduct.IsEmpty() {
		return errors.Wrapf(error_codes.ErrProductNotExist, "postgres.PriceTimeStampAdd")
	}

	query, args, err := squirrel.Insert("price_history").
		Columns("code, price, date").
		Values(code, priceTimeStamp.Price, priceTimeStamp.Date).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return errors.Wrapf(err, "postgres.PriceTimeStampAdd: to sql")
	}

	_, err = p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "postgres.PriceTimeStampAdd: to query")
	}

	return nil
}

func (p *postgres) PriceHistory(ctx context.Context, code string) (models.PriceHistory, error) {
	existingProduct, err := p.ProductGet(ctx, code)
	if err != nil {
		return nil, err
	}
	if existingProduct.IsEmpty() {
		return nil, errors.Wrapf(error_codes.ErrProductNotExist, "postgres.PriceHistory")
	}

	query, args, err := squirrel.Select("price, date").
		From("price_history").
		Where(squirrel.Eq{
			"code": code,
		}).
		OrderBy("date").
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return nil, errors.Wrapf(err, "postgres.PriceHistory: to sql")
	}

	var priceHistory models.PriceHistory
	err = pgxscan.Select(ctx, p.pool, &priceHistory, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "postgres.PriceHistory: to query")
	}

	return priceHistory, nil
}
