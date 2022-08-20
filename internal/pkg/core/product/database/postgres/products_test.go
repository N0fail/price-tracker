package postgres

import (
	"context"
	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"regexp"
	"testing"
	"time"
)

func TestProductGet(t *testing.T) {
	t.Run("get no products", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs("123").
			WillReturnRows(rows)

		// act
		product, err := f.postgres.ProductGet(ctx, "123")

		// assert
		require.NoError(t, err)
		require.Equal(t, true, product.IsEmpty())
	})

	t.Run("get product", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow("123", "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs("123").
			WillReturnRows(rows)

		// act
		product, err := f.postgres.ProductGet(ctx, "123")

		// assert
		require.NoError(t, err)
		require.Equal(t, "123", product.Code)
		require.Equal(t, "name", product.Name)
	})

	t.Run("get product multiple products", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow("123", "name").AddRow("123", "name2")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs("123").
			WillReturnRows(rows)

		// act
		_, err := f.postgres.ProductGet(ctx, "123")

		// assert
		require.Equal(t, "postgres.ProductGet: found mutiple producs with code 123", err.Error())
	})

	t.Run("get product database error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs("123").
			WillReturnError(ErrQuery)

		// act
		_, err := f.postgres.ProductGet(ctx, "123")

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})
}

func TestProductCreate(t *testing.T) {
	t.Run("get product error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode, productName := "123", "name"
		product := models.Product{
			Name: productName,
			Code: productCode,
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnError(ErrQuery)

		// act
		err := f.postgres.ProductCreate(ctx, product)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("already exist", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode, productName := "123", "name"
		product := models.Product{
			Name: productName,
			Code: productCode,
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow("123", "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		// act
		err := f.postgres.ProductCreate(ctx, product)

		// assert
		require.Equal(t, true, errors.Is(err, error_codes.ErrProductExists))
	})

	t.Run("error in insert query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode, productName := "123", "name"
		product := models.Product{
			Name: productName,
			Code: productCode,
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)
		querySql = regexp.QuoteMeta(`INSERT INTO products (code, name) VALUES ($1,$2)`)
		f.dbConnMock.ExpectExec(querySql).
			WithArgs(productCode, productName).
			WillReturnError(ErrQuery)

		// act
		err := f.postgres.ProductCreate(ctx, product)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("success insert query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode, productName := "123", "name"
		product := models.Product{
			Name: productName,
			Code: productCode,
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)
		querySql = regexp.QuoteMeta(`INSERT INTO products (code, name) VALUES ($1,$2)`)
		f.dbConnMock.ExpectExec(querySql).
			WithArgs(productCode, productName).
			WillReturnResult(pgxmock.NewResult("", 1))

		// act
		err := f.postgres.ProductCreate(ctx, product)

		// assert
		require.Nil(t, err)
	})
}

func TestProductDelete(t *testing.T) {
	t.Run("get product error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnError(ErrQuery)

		// act
		err := f.postgres.ProductDelete(ctx, productCode)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("product not exist", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		// act
		err := f.postgres.ProductDelete(ctx, productCode)

		// assert
		require.Equal(t, true, errors.Is(err, error_codes.ErrProductNotExist))
	})

	t.Run("error in delete query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow(productCode, "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		querySql = regexp.QuoteMeta("DELETE FROM products WHERE code = $1")
		f.dbConnMock.ExpectExec(querySql).
			WithArgs(productCode).
			WillReturnError(ErrQuery)

		// act
		err := f.postgres.ProductDelete(ctx, productCode)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("successful delete query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow(productCode, "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		querySql = regexp.QuoteMeta("DELETE FROM products WHERE code = $1")
		f.dbConnMock.ExpectExec(querySql).
			WithArgs(productCode).
			WillReturnResult(pgxmock.NewResult("", 1))

		// act
		err := f.postgres.ProductDelete(ctx, productCode)

		// assert
		require.Nil(t, err)
	})
}

func TestPriceTimeStampAdd(t *testing.T) {
	t.Run("get product error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		priceTimeStamp := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnError(ErrQuery)

		// act
		err := f.postgres.PriceTimeStampAdd(ctx, productCode, priceTimeStamp)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("product not exist", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		priceTimeStamp := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		// act
		err := f.postgres.PriceTimeStampAdd(ctx, productCode, priceTimeStamp)

		// assert
		require.Equal(t, true, errors.Is(err, error_codes.ErrProductNotExist))
	})

	t.Run("error in insert query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		priceTimeStamp := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow(productCode, "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		querySql = regexp.QuoteMeta("INSERT INTO price_history (code, price, date) VALUES ($1,$2,$3)")
		f.dbConnMock.ExpectExec(querySql).
			WithArgs(productCode, priceTimeStamp.Price, priceTimeStamp.Date).
			WillReturnError(ErrQuery)

		// act
		err := f.postgres.PriceTimeStampAdd(ctx, productCode, priceTimeStamp)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("successful insert query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		priceTimeStamp := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		}
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow(productCode, "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		querySql = regexp.QuoteMeta("INSERT INTO price_history (code, price, date) VALUES ($1,$2,$3)")
		f.dbConnMock.ExpectExec(querySql).
			WithArgs(productCode, priceTimeStamp.Price, priceTimeStamp.Date).
			WillReturnResult(pgxmock.NewResult("", 1))

		// act
		err := f.postgres.PriceTimeStampAdd(ctx, productCode, priceTimeStamp)

		// assert
		require.Nil(t, err)
	})
}

func TestFullHistory(t *testing.T) {
	t.Run("get product error", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnError(ErrQuery)

		// act
		_, err := f.postgres.PriceHistory(ctx, productCode)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("product not exist", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		// act
		_, err := f.postgres.PriceHistory(ctx, productCode)

		// assert
		require.Equal(t, true, errors.Is(err, error_codes.ErrProductNotExist))
	})

	t.Run("error in history query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		productCode := "123"
		querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
		rows := pgxmock.NewRows([]string{"code", "name"}).AddRow(productCode, "name")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnRows(rows)

		querySql = regexp.QuoteMeta("SELECT price, date FROM price_history WHERE code = $1 ORDER BY date")
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(productCode).
			WillReturnError(ErrQuery)

		// act
		_, err := f.postgres.PriceHistory(ctx, productCode)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
	})

	t.Run("successful history queries", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)
		for _, historyVariant := range []models.PriceHistory{
			{models.PriceTimeStamp{Price: 1, Date: time.Unix(1, 1)}},
			{models.PriceTimeStamp{Price: 1, Date: time.Unix(1, 1)}, models.PriceTimeStamp{Price: 2, Date: time.Unix(2, 2)}},
			{},
		} {
			productCode := "123"
			querySql := regexp.QuoteMeta("SELECT code, name FROM products WHERE code = $1")
			rows := pgxmock.NewRows([]string{"code", "name"}).AddRow(productCode, "name")
			f.dbConnMock.ExpectQuery(querySql).
				WithArgs(productCode).
				WillReturnRows(rows)

			querySql = regexp.QuoteMeta("SELECT price, date FROM price_history WHERE code = $1 ORDER BY date")
			historyRows := pgxmock.NewRows([]string{"price", "date"})
			for _, priceTimeStamp := range historyVariant {
				historyRows = historyRows.AddRow(priceTimeStamp.Price, priceTimeStamp.Date)
			}

			f.dbConnMock.ExpectQuery(querySql).
				WithArgs(productCode).
				WillReturnRows(historyRows)

			// act
			history, err := f.postgres.PriceHistory(ctx, productCode)

			// assert
			require.Nil(t, err)
			require.Equal(t, len(historyVariant), len(history))
			for i, priceTimeStamp := range historyVariant {
				require.Equal(t, priceTimeStamp.Price, history[i].Price)
				require.Equal(t, priceTimeStamp.Date, history[i].Date)
			}
		}
	})
}

func TestProductList(t *testing.T) {
	t.Run("error in query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta(`
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
		`)
		orderBy := "code"
		pageNumber, resultsPerPage := 1, 10
		limit, offset := resultsPerPage, pageNumber*resultsPerPage

		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(uint32(limit), uint32(offset), orderBy).
			WillReturnError(ErrQuery)

		// act
		list, err := f.postgres.ProductList(ctx, uint32(pageNumber), uint32(resultsPerPage), orderBy)

		// assert
		require.Equal(t, true, errors.Is(err, ErrQuery))
		require.Nil(t, list)
	})

	t.Run("error no entries", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta(`
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
		`)
		orderBy := "code"
		pageNumber, resultsPerPage := 1, 10
		limit, offset := resultsPerPage, pageNumber*resultsPerPage
		rows := pgxmock.NewRows([]string{"products.code", "products.name", "last_price.price", "last_price.date"})
		f.dbConnMock.ExpectQuery(querySql).
			WithArgs(uint32(limit), uint32(offset), orderBy).
			WillReturnRows(rows)

		// act
		list, err := f.postgres.ProductList(ctx, uint32(pageNumber), uint32(resultsPerPage), orderBy)

		// assert
		require.Equal(t, true, errors.Is(err, error_codes.ErrNoEntries))
		require.Nil(t, list)
	})

	t.Run("error in query", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		f := setUp(t)
		defer f.tearDown(ctx)

		querySql := regexp.QuoteMeta(`
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
		`)
		orderBy := "code"
		pageNumber, resultsPerPage := 1, 10
		limit, offset := resultsPerPage, pageNumber*resultsPerPage
		for _, listVariant := range [][]models.ProductSnapshot{
			{
				models.ProductSnapshot{Code: "1", Name: "1n", LastPrice: models.PriceTimeStamp{Price: 1, Date: time.Unix(1, 1)}},
				models.ProductSnapshot{Code: "2", Name: "2n", LastPrice: models.EmptyPriceTimeStamp},
			},
			{
				models.ProductSnapshot{Code: "2", Name: "2n", LastPrice: models.EmptyPriceTimeStamp},
			},
			{
				models.ProductSnapshot{Code: "1", Name: "1n", LastPrice: models.PriceTimeStamp{Price: 1, Date: time.Unix(1, 1)}},
			},
		} {
			rows := pgxmock.NewRows([]string{"products.code", "products.name", "last_price.price", "last_price.date"})
			for _, productSnapshot := range listVariant {
				if productSnapshot.LastPrice.IsEmpty() {
					rows = rows.AddRow(productSnapshot.Code, productSnapshot.Name, nil, nil)
				} else {
					rows = rows.AddRow(productSnapshot.Code, productSnapshot.Name, productSnapshot.LastPrice.Price, productSnapshot.LastPrice.Date)
				}
			}
			f.dbConnMock.ExpectQuery(querySql).
				WithArgs(uint32(limit), uint32(offset), orderBy).
				WillReturnRows(rows)

			// act
			list, err := f.postgres.ProductList(ctx, uint32(pageNumber), uint32(resultsPerPage), orderBy)

			// assert
			require.Nil(t, err)
			require.Equal(t, len(listVariant), len(list))
			for i, productSnapshot := range listVariant {
				require.Equal(t, productSnapshot.Name, list[i].Name)
				require.Equal(t, productSnapshot.Code, list[i].Code)
				if productSnapshot.LastPrice.IsEmpty() {
					require.Equal(t, true, list[i].LastPrice.IsEmpty())
				} else {
					require.Equal(t, productSnapshot.LastPrice.Price, list[i].LastPrice.Price)
					require.Equal(t, productSnapshot.LastPrice.Date, list[i].LastPrice.Date)
				}
			}
		}
	})
}
