package api

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestProductCreateApi(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		p := productSetUp(t)

		productCode := "123"
		productName := "qwerty"
		productModel := models.Product{
			Code: productCode,
			Name: productName,
		}

		p.product.EXPECT().
			ProductCreate(gomock.Any(), productModel).Return(nil)

		// act
		resp, err := p.api.ProductCreate(context.Background(), &pb.ProductCreateRequest{
			Code: productCode,
			Name: productName,
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, resp, &pb.ProductCreateResponse{})
	})

	t.Run("error product exists", func(t *testing.T) {
		// arrange
		p := productSetUp(t)

		productCode := "123"
		productName := "qwerty"
		productModel := models.Product{
			Code: productCode,
			Name: productName,
		}

		p.product.EXPECT().
			ProductCreate(gomock.Any(), productModel).Return(error_codes.ErrProductExists)

		// act
		resp, err := p.api.ProductCreate(context.Background(), &pb.ProductCreateRequest{
			Code: productCode,
			Name: productName,
		})

		// assert
		require.Nil(t, resp)
		require.ErrorIs(t, err, status.Error(codes.AlreadyExists, error_codes.ErrProductExists.Error()))
	})

	t.Run("error internal", func(t *testing.T) {
		// arrange
		p := productSetUp(t)

		productCode := "123"
		productName := "qwerty"
		productModel := models.Product{
			Code: productCode,
			Name: productName,
		}

		p.product.EXPECT().
			ProductCreate(gomock.Any(), productModel).Return(error_codes.ErrExternalProblem)

		// act
		resp, err := p.api.ProductCreate(context.Background(), &pb.ProductCreateRequest{
			Code: productCode,
			Name: productName,
		})

		// assert
		require.Nil(t, resp)
		require.ErrorIs(t, err, status.Error(codes.Internal, error_codes.ErrExternalProblem.Error()))
	})
}

func TestProductListApi(t *testing.T) {
	t.Run("success order by code", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		pageNumber, resultsPerPage, orderBy := uint32(1), uint32(10), "code"
		productCode, productName := "1234", "qwerty"
		date := time.Unix(123, 0)
		productPrice := models.PriceTimeStamp{Price: 322, Date: date}

		p.product.EXPECT().
			ProductList(gomock.Any(), pageNumber, resultsPerPage, orderBy).Return([]models.ProductSnapshot{{
			Code:      productCode,
			Name:      productName,
			LastPrice: productPrice,
		}}, nil)

		// act
		resp, err := p.api.ProductList(context.Background(), &pb.ProductListRequest{
			PageNumber:     pageNumber,
			ResultsPerPage: resultsPerPage,
			OrderBy:        pb.ProductListRequest_OrderBy(pb.ProductListRequest_OrderBy_value[orderBy]),
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, resp, &pb.ProductListResponse{
			ProductSnapShots: []*pb.ProductSnapShot{{
				Code: productCode,
				Name: productName,
				PriceTimeStamp: &pb.PriceTimeStamp{
					Price: productPrice.Price,
					Ts:    date.Unix(),
				},
			}},
		})
	})

	t.Run("success order by name", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		pageNumber, resultsPerPage, orderBy := uint32(3), uint32(17), "name"
		productCode, productName := "789", "asdf"
		date := time.Unix(456, 0)
		productPrice := models.PriceTimeStamp{Price: 777, Date: date}

		p.product.EXPECT().
			ProductList(gomock.Any(), pageNumber, resultsPerPage, orderBy).Return([]models.ProductSnapshot{{
			Code:      productCode,
			Name:      productName,
			LastPrice: productPrice,
		}}, nil)

		// act
		resp, err := p.api.ProductList(context.Background(), &pb.ProductListRequest{
			PageNumber:     pageNumber,
			ResultsPerPage: resultsPerPage,
			OrderBy:        pb.ProductListRequest_OrderBy(pb.ProductListRequest_OrderBy_value[orderBy]),
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, resp, &pb.ProductListResponse{
			ProductSnapShots: []*pb.ProductSnapShot{{
				Code: productCode,
				Name: productName,
				PriceTimeStamp: &pb.PriceTimeStamp{
					Price: productPrice.Price,
					Ts:    date.Unix(),
				},
			}},
		})
	})

	t.Run("success no LastPrice", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		pageNumber, resultsPerPage, orderBy := uint32(3), uint32(17), "name"
		productCode, productName := "789", "asdf"

		p.product.EXPECT().
			ProductList(gomock.Any(), pageNumber, resultsPerPage, orderBy).Return([]models.ProductSnapshot{{
			Code:      productCode,
			Name:      productName,
			LastPrice: models.EmptyPriceTimeStamp,
		}}, nil)

		// act
		resp, err := p.api.ProductList(context.Background(), &pb.ProductListRequest{
			PageNumber:     pageNumber,
			ResultsPerPage: resultsPerPage,
			OrderBy:        pb.ProductListRequest_OrderBy(pb.ProductListRequest_OrderBy_value[orderBy]),
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, resp, &pb.ProductListResponse{
			ProductSnapShots: []*pb.ProductSnapShot{{
				Code:           productCode,
				Name:           productName,
				PriceTimeStamp: nil,
			}},
		})
	})

	t.Run("error no entries", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		pageNumber, resultsPerPage, orderBy := uint32(3), uint32(17), "name"

		p.product.EXPECT().
			ProductList(gomock.Any(), pageNumber, resultsPerPage, orderBy).Return(nil, error_codes.ErrNoEntries)

		// act
		resp, err := p.api.ProductList(context.Background(), &pb.ProductListRequest{
			PageNumber:     pageNumber,
			ResultsPerPage: resultsPerPage,
			OrderBy:        pb.ProductListRequest_OrderBy(pb.ProductListRequest_OrderBy_value[orderBy]),
		})

		// assert
		require.ErrorIs(t, err, error_codes.ErrNoEntries)
		require.Nil(t, resp)
	})

	t.Run("error external", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		pageNumber, resultsPerPage, orderBy := uint32(3), uint32(17), "name"

		p.product.EXPECT().
			ProductList(gomock.Any(), pageNumber, resultsPerPage, orderBy).Return(nil, errors.New("Some error"))

		// act
		resp, err := p.api.ProductList(context.Background(), &pb.ProductListRequest{
			PageNumber:     pageNumber,
			ResultsPerPage: resultsPerPage,
			OrderBy:        pb.ProductListRequest_OrderBy(pb.ProductListRequest_OrderBy_value[orderBy]),
		})

		// assert
		require.ErrorIs(t, err, error_codes.ErrExternalProblem)
		require.Nil(t, resp)
	})
}

func TestProductDeleteApi(t *testing.T) {
	t.Run("success delete", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"

		p.product.EXPECT().
			ProductDelete(gomock.Any(), productCode).Return(nil)

		// act
		resp, err := p.api.ProductDelete(context.Background(), &pb.ProductDeleteRequest{
			Code: productCode,
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, &pb.ProductDeleteResponse{}, resp)
	})

	t.Run("error does not exist", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"

		p.product.EXPECT().
			ProductDelete(gomock.Any(), productCode).Return(error_codes.ErrProductNotExist)

		// act
		resp, err := p.api.ProductDelete(context.Background(), &pb.ProductDeleteRequest{
			Code: productCode,
		})

		// assert
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, error_codes.ErrProductNotExist.Error()))
		require.Nil(t, resp)
	})

	t.Run("error internal", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"

		p.product.EXPECT().
			ProductDelete(gomock.Any(), productCode).Return(error_codes.ErrExternalProblem)

		// act
		resp, err := p.api.ProductDelete(context.Background(), &pb.ProductDeleteRequest{
			Code: productCode,
		})

		// assert
		require.ErrorIs(t, err, status.Error(codes.Internal, error_codes.ErrExternalProblem.Error()))
		require.Nil(t, resp)
	})
}

func TestPriceTimeStampAddApi(t *testing.T) {
	t.Run("success add", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"
		price := float64(322)
		date := time.Unix(123, 0)
		priceTimeStamp := models.PriceTimeStamp{Price: price, Date: date}

		p.product.EXPECT().
			PriceTimeStampAdd(gomock.Any(), productCode, priceTimeStamp).Return(nil)

		// act
		resp, err := p.api.PriceTimeStampAdd(context.Background(), &pb.PriceTimeStampAddRequest{
			Code:  productCode,
			Ts:    date.Unix(),
			Price: price,
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, &pb.PriceTimeStampAddResponse{}, resp)
	})

	t.Run("error product not exist", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"
		price := float64(322)
		date := time.Unix(123, 0)
		priceTimeStamp := models.PriceTimeStamp{Price: price, Date: date}

		p.product.EXPECT().
			PriceTimeStampAdd(gomock.Any(), productCode, priceTimeStamp).Return(error_codes.ErrProductNotExist)

		// act
		resp, err := p.api.PriceTimeStampAdd(context.Background(), &pb.PriceTimeStampAddRequest{
			Code:  productCode,
			Ts:    date.Unix(),
			Price: price,
		})

		// assert
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, error_codes.ErrProductNotExist.Error()))
		require.Nil(t, resp)
	})

	t.Run("error internal", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"
		price := float64(322)
		date := time.Unix(123, 0)
		priceTimeStamp := models.PriceTimeStamp{Price: price, Date: date}

		p.product.EXPECT().
			PriceTimeStampAdd(gomock.Any(), productCode, priceTimeStamp).Return(error_codes.ErrExternalProblem)

		// act
		resp, err := p.api.PriceTimeStampAdd(context.Background(), &pb.PriceTimeStampAddRequest{
			Code:  productCode,
			Ts:    date.Unix(),
			Price: price,
		})

		// assert
		require.ErrorIs(t, err, status.Error(codes.Internal, error_codes.ErrExternalProblem.Error()))
		require.Nil(t, resp)
	})
}

func TestPriceHistoryApi(t *testing.T) {
	t.Run("success history", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"
		price := float64(322)
		date := time.Unix(123, 0)
		priceTimeStamp := models.PriceTimeStamp{Price: price, Date: date}

		p.product.EXPECT().
			PriceHistory(gomock.Any(), productCode).Return(models.PriceHistory{priceTimeStamp}, nil)

		// act
		resp, err := p.api.PriceHistory(context.Background(), &pb.PriceHistoryRequest{
			Code: productCode,
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, &pb.PriceHistoryResponse{
			PriceHistory: []*pb.PriceTimeStamp{{
				Price: price,
				Ts:    date.Unix(),
			}},
		}, resp)
	})

	t.Run("success history 2 entries", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"
		price1, price2 := float64(322), float64(876)
		date1, date2 := time.Unix(123, 0), time.Unix(765, 0)
		priceTimeStamp1 := models.PriceTimeStamp{Price: price1, Date: date1}
		priceTimeStamp2 := models.PriceTimeStamp{Price: price2, Date: date2}

		p.product.EXPECT().
			PriceHistory(gomock.Any(), productCode).Return(models.PriceHistory{priceTimeStamp1, priceTimeStamp2}, nil)

		// act
		resp, err := p.api.PriceHistory(context.Background(), &pb.PriceHistoryRequest{
			Code: productCode,
		})

		// assert
		require.NoError(t, err)
		require.Equal(t, &pb.PriceHistoryResponse{
			PriceHistory: []*pb.PriceTimeStamp{
				{
					Price: price1,
					Ts:    date1.Unix(),
				},
				{
					Price: price2,
					Ts:    date2.Unix(),
				},
			},
		}, resp)
	})

	t.Run("error product not exist", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"

		p.product.EXPECT().
			PriceHistory(gomock.Any(), productCode).Return(nil, error_codes.ErrProductNotExist)

		// act
		resp, err := p.api.PriceHistory(context.Background(), &pb.PriceHistoryRequest{
			Code: productCode,
		})

		// assert
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, error_codes.ErrProductNotExist.Error()))
		require.Nil(t, resp)
	})

	t.Run("error internal", func(t *testing.T) {
		// arrange
		p := productSetUp(t)
		productCode := "1234"

		p.product.EXPECT().
			PriceHistory(gomock.Any(), productCode).Return(nil, error_codes.ErrExternalProblem)

		// act
		resp, err := p.api.PriceHistory(context.Background(), &pb.PriceHistoryRequest{
			Code: productCode,
		})

		// assert
		require.ErrorIs(t, err, status.Error(codes.Internal, error_codes.ErrExternalProblem.Error()))
		require.Nil(t, resp)
	})
}
