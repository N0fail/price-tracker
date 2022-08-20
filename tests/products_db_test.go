//go:build integration
// +build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	"testing"
	"time"
)

func TestProductCreateDb(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		err := productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})

		//assert
		require.NoError(t, err)
	})

	t.Run("error create twice", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		err1 := productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		err2 := productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})

		//assert
		require.NoError(t, err1)
		require.ErrorIs(t, err2, error_codes.ErrProductExists)
	})
}

func TestProductDeleteDb(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		err := productCore.ProductDelete(context.Background(), "123")

		//assert
		require.NoError(t, err)
	})

	t.Run("error create twice", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		err := productCore.ProductDelete(context.Background(), "123")

		//assert
		require.ErrorIs(t, err, error_codes.ErrProductNotExist)
	})
}

func TestPriceTimeStampAddDb(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		err := productCore.PriceTimeStampAdd(context.Background(), "123", models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		})

		//assert
		require.NoError(t, err)
	})

	t.Run("success twice", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		err1 := productCore.PriceTimeStampAdd(context.Background(), "123", models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		})
		err2 := productCore.PriceTimeStampAdd(context.Background(), "123", models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(124, 0),
		})

		//assert
		require.NoError(t, err1)
		require.NoError(t, err2)
	})

	t.Run("no product", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		err := productCore.PriceTimeStampAdd(context.Background(), "123", models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0),
		})

		//assert
		require.ErrorIs(t, err, error_codes.ErrProductNotExist)
	})
}

func TestPriceHistoryDb(t *testing.T) {
	t.Run("success no history", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		history, err := productCore.PriceHistory(context.Background(), "123")

		//assert
		require.NoError(t, err)
		require.Nil(t, history)
	})

	t.Run("success 1 entry", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		priceTimeStamp := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0).UTC(),
		}
		productCore.PriceTimeStampAdd(context.Background(), "123", priceTimeStamp)
		history, err := productCore.PriceHistory(context.Background(), "123")

		//assert
		require.NoError(t, err)
		require.Equal(t, models.PriceHistory{priceTimeStamp}, history)
	})

	t.Run("success 2 entries sorted", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		productCore.ProductCreate(context.Background(), models.Product{
			Code: "123",
			Name: "qwerty",
		})
		priceTimeStamp2 := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0).UTC(),
		}
		priceTimeStamp1 := models.PriceTimeStamp{
			Price: 123,
			Date:  time.Unix(23, 0).UTC(),
		}
		productCore.PriceTimeStampAdd(context.Background(), "123", priceTimeStamp2)
		productCore.PriceTimeStampAdd(context.Background(), "123", priceTimeStamp1)
		history, err := productCore.PriceHistory(context.Background(), "123")

		//assert
		require.NoError(t, err)
		require.Equal(t, models.PriceHistory{priceTimeStamp1, priceTimeStamp2}, history)
	})

	t.Run("error no product", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		history, err := productCore.PriceHistory(context.Background(), "123")

		//assert
		require.ErrorIs(t, err, error_codes.ErrProductNotExist)
		require.Nil(t, history)
	})
}

func TestProductListDb(t *testing.T) {
	t.Run("success 1 entry no price", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		product1 := models.Product{
			Code: "123",
			Name: "qwerty",
		}
		productSnapshot1 := models.ProductSnapshot{
			Code:      product1.Code,
			Name:      product1.Name,
			LastPrice: models.EmptyPriceTimeStamp,
		}
		productCore.ProductCreate(context.Background(), product1)
		pageNumber, resultsPerPage, orderBy := uint32(0), uint32(2), "code"
		list, err := productCore.ProductList(context.Background(), pageNumber, resultsPerPage, orderBy)

		//assert
		require.NoError(t, err)
		require.Equal(t, []models.ProductSnapshot{productSnapshot1}, list)
	})

	t.Run("success 2 entries: no price + price", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()
		product1 := models.Product{
			Code: "123",
			Name: "qwerty",
		}
		price1 := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0).UTC(),
		}
		price2 := models.PriceTimeStamp{
			Price: 555,
			Date:  time.Unix(1, 0).UTC(),
		}
		productSnapshot1 := models.ProductSnapshot{
			Code:      product1.Code,
			Name:      product1.Name,
			LastPrice: price1,
		}
		product2 := models.Product{
			Code: "456",
			Name: "asdfg",
		}
		productSnapshot2 := models.ProductSnapshot{
			Code:      product2.Code,
			Name:      product2.Name,
			LastPrice: models.EmptyPriceTimeStamp,
		}

		//act
		productCore.ProductCreate(context.Background(), product1)
		productCore.ProductCreate(context.Background(), product2)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price1)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price2)
		pageNumber, resultsPerPage, orderBy := uint32(0), uint32(2), "code"
		list, err := productCore.ProductList(context.Background(), pageNumber, resultsPerPage, orderBy)

		//assert
		require.NoError(t, err)
		require.Equal(t, []models.ProductSnapshot{productSnapshot1, productSnapshot2}, list)
	})

	t.Run("success 1 entry order by code", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()
		product1 := models.Product{
			Code: "123",
			Name: "qwerty",
		}
		price1 := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0).UTC(),
		}
		price2 := models.PriceTimeStamp{
			Price: 555,
			Date:  time.Unix(1, 0).UTC(),
		}
		productSnapshot1 := models.ProductSnapshot{
			Code:      product1.Code,
			Name:      product1.Name,
			LastPrice: price1,
		}
		product2 := models.Product{
			Code: "456",
			Name: "asdfg",
		}
		//productSnapshot2 := models.ProductSnapshot{
		//	Code:      product2.Code,
		//	Name:      product2.Name,
		//	LastPrice: models.EmptyPriceTimeStamp,
		//}

		//act
		productCore.ProductCreate(context.Background(), product1)
		productCore.ProductCreate(context.Background(), product2)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price1)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price2)
		pageNumber, resultsPerPage, orderBy := uint32(0), uint32(1), "code"
		list, err := productCore.ProductList(context.Background(), pageNumber, resultsPerPage, orderBy)

		//assert
		require.NoError(t, err)
		require.Equal(t, []models.ProductSnapshot{productSnapshot1}, list)
	})

	t.Run("success 1 entry order by name", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()
		product1 := models.Product{
			Code: "123",
			Name: "qwerty",
		}
		price1 := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0).UTC(),
		}
		price2 := models.PriceTimeStamp{
			Price: 555,
			Date:  time.Unix(1, 0).UTC(),
		}
		//productSnapshot1 := models.ProductSnapshot{
		//	Code:      product1.Code,
		//	Name:      product1.Name,
		//	LastPrice: price1,
		//}
		product2 := models.Product{
			Code: "456",
			Name: "asdfg",
		}
		productSnapshot2 := models.ProductSnapshot{
			Code:      product2.Code,
			Name:      product2.Name,
			LastPrice: models.EmptyPriceTimeStamp,
		}

		//act
		productCore.ProductCreate(context.Background(), product1)
		productCore.ProductCreate(context.Background(), product2)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price1)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price2)
		pageNumber, resultsPerPage, orderBy := uint32(0), uint32(1), "name"
		list, err := productCore.ProductList(context.Background(), pageNumber, resultsPerPage, orderBy)

		//assert
		require.NoError(t, err)
		require.Equal(t, []models.ProductSnapshot{productSnapshot2}, list)
	})

	t.Run("error no entries in given range", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()
		product1 := models.Product{
			Code: "123",
			Name: "qwerty",
		}
		price1 := models.PriceTimeStamp{
			Price: 322,
			Date:  time.Unix(123, 0).UTC(),
		}
		price2 := models.PriceTimeStamp{
			Price: 555,
			Date:  time.Unix(1, 0).UTC(),
		}
		//productSnapshot1 := models.ProductSnapshot{
		//	Code:      product1.Code,
		//	Name:      product1.Name,
		//	LastPrice: price1,
		//}
		product2 := models.Product{
			Code: "456",
			Name: "asdfg",
		}
		//productSnapshot2 := models.ProductSnapshot{
		//	Code:      product2.Code,
		//	Name:      product2.Name,
		//	LastPrice: models.EmptyPriceTimeStamp,
		//}

		//act
		productCore.ProductCreate(context.Background(), product1)
		productCore.ProductCreate(context.Background(), product2)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price1)
		productCore.PriceTimeStampAdd(context.Background(), product1.Code, price2)
		pageNumber, resultsPerPage, orderBy := uint32(2), uint32(1), "name"
		list, err := productCore.ProductList(context.Background(), pageNumber, resultsPerPage, orderBy)

		//assert
		require.ErrorIs(t, err, error_codes.ErrNoEntries)
		require.Nil(t, list)
	})

	t.Run("error no entries", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		productCore := product.New(Db.Pool)
		defer Db.TearDown()

		//act
		pageNumber, resultsPerPage, orderBy := uint32(0), uint32(1), "name"
		list, err := productCore.ProductList(context.Background(), pageNumber, resultsPerPage, orderBy)

		//assert
		require.ErrorIs(t, err, error_codes.ErrNoEntries)
		require.Nil(t, list)
	})
}
