package api

import (
	"context"
	"github.com/pkg/errors"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func New(p productPkg.Interface) pb.AdminServer {
	return &implementation{
		product: p,
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	product productPkg.Interface
}

func (i *implementation) ProductCreate(ctx context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	if err := i.product.Create(models.Product{
		Code: in.GetCode(),
		Name: in.GetName(),
	}); err != nil {
		if errors.Is(err, error_codes.ErrNameTooShortError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductCreateResponse{}, nil
}
func (i *implementation) ProductList(context.Context, *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	productSnapShots := i.product.List()

	result := make([]*pb.ProductSnapShot, 0, len(productSnapShots))
	for _, productSnapShot := range productSnapShots {
		result = append(result, &pb.ProductSnapShot{
			Code: productSnapShot.Code,
			Name: productSnapShot.Name,
			PriceTimeStamp: &pb.PriceTimeStamp{
				Price: productSnapShot.LastPrice.Price,
				Ts:    productSnapShot.LastPrice.Date.Unix(),
			},
		})
	}

	return &pb.ProductListResponse{
		ProductSnapShots: result,
	}, nil
}
func (i *implementation) ProductDelete(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	if err := i.product.Delete(in.Code); err != nil {
		if errors.Is(err, error_codes.ErrProductNotExist) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductDeleteResponse{}, nil
}
func (i *implementation) PriceTimeStampAdd(ctx context.Context, in *pb.PriceTimeStampAddRequest) (*pb.PriceTimeStampAddResponse, error) {
	ts := models.PriceTimeStamp{
		Price: in.GetPrice(),
		Date:  time.Unix(in.GetTs(), 0),
	}
	if err := i.product.AddPriceTimeStamp(in.Code, ts); err != nil {
		if errors.Is(err, error_codes.ErrProductNotExist) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PriceTimeStampAddResponse{}, nil
}
func (i *implementation) PriceHistory(ctx context.Context, in *pb.PriceHistoryRequest) (*pb.PriceHistoryResponse, error) {
	priceHistory, err := i.product.FullHistory(in.GetCode())
	if err != nil {
		if errors.Is(err, error_codes.ErrProductNotExist) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.PriceTimeStamp, 0, len(priceHistory))
	for _, priceTimeStamp := range priceHistory {
		result = append(result, &pb.PriceTimeStamp{
			Price: priceTimeStamp.Price,
			Ts:    priceTimeStamp.Date.Unix(),
		})
	}

	return &pb.PriceHistoryResponse{
		PriceHistory: result,
	}, nil
}
