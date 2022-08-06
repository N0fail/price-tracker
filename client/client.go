package main

import (
	"context"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"log"
	"time"

	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// у клиента поидее нет доступа к internal константам, но для удобства можно написать
	conns, err := grpc.Dial(config.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAdminClient(conns)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	list(client, ctx, "empty")

	productCreate(client, ctx, "1", "333", "name too short")
	productCreate(client, ctx, "1", "4444", "ok")
	productCreate(client, ctx, "1", "4444", "already exists")
	productCreate(client, ctx, "2", "5555", "ok")

	list(client, ctx, "2 entries")

	priceTimeStampAdd(client, ctx, "3", "2 Jan 2006 15:04", 23.45, "product dont exist")
	priceTimeStampAdd(client, ctx, "2", "2 Jan 2006 15:04", -23.45, "negative price")
	priceTimeStampAdd(client, ctx, "2", "2 Jan 2006 15:04", 23.45, "ok")
	priceTimeStampAdd(client, ctx, "2", "1 Jan 2006 15:04", 34.45, "ok")

	list(client, ctx, "2 entries, 1 with price")
	priceHistory(client, ctx, "2", "2 records")
	priceHistory(client, ctx, "1", "empty")
	priceHistory(client, ctx, "3", "product dont exist")

	delete(client, ctx, "3", "product dont exist")
	delete(client, ctx, "1", "ok")
	delete(client, ctx, "1", "product dont exist")

	list(client, ctx, "1 entry with price")
	priceHistory(client, ctx, "1", "product dont exist")

	delete(client, ctx, "2", "ok")
	list(client, ctx, "empty")
}

func list(client pb.AdminClient, ctx context.Context, expected string) {
	funcName := "list"
	resp, err := client.ProductList(ctx, &pb.ProductListRequest{
		Page: 0,
	})
	if err != nil {
		log.Printf("func: [%v]:\nerror: [%v]\nexpected: [%v]", funcName, err, expected)
		return
	}

	log.Printf("func: [%v]:\nresponse: [%v]\nexpected: [%v]", funcName, resp, expected)
}

func productCreate(client pb.AdminClient, ctx context.Context, code, name string, expected string) {
	funcName := "productCreate"
	resp, err := client.ProductCreate(ctx, &pb.ProductCreateRequest{
		Code: code,
		Name: name,
	})
	if err != nil {
		log.Printf("func: [%v]:\nerror: [%v]\nexpected: [%v]", funcName, err, expected)
		return
	}

	log.Printf("func: [%v]:\nresponse: [%v]\nexpected: [%v]", funcName, resp, expected)
}

func priceTimeStampAdd(client pb.AdminClient, ctx context.Context, code, date string, price float64, expected string) {
	funcName := "priceTimeStampAdd"
	priceTime, err := time.Parse("2 Jan 2006 15:04", date)
	if err != nil {
		log.Printf("func: [%v]:\nerror: [%v]\nexpected: [%v]", funcName, err, expected)
		return
	}
	resp, err := client.PriceTimeStampAdd(ctx, &pb.PriceTimeStampAddRequest{
		Price: price,
		Code:  code,
		Ts:    priceTime.Unix(),
	})

	if err != nil {
		log.Printf("func: [%v]:\nerror: [%v]\nexpected: [%v]", funcName, err, expected)
		return
	}

	log.Printf("func: [%v]:\nresponse: [%v]\nexpected: [%v]", funcName, resp, expected)
}

func priceHistory(client pb.AdminClient, ctx context.Context, code string, expected string) {
	funcName := "priceHistory"
	resp, err := client.PriceHistory(ctx, &pb.PriceHistoryRequest{
		Code: code,
	})

	if err != nil {
		log.Printf("func: [%v]:\nerror: [%v]\nexpected: [%v]", funcName, err, expected)
		return
	}

	log.Printf("func: [%v]:\nresponse: [%v]\nexpected: [%v]", funcName, resp, expected)
}

func delete(client pb.AdminClient, ctx context.Context, code string, expected string) {
	funcName := "delete"
	resp, err := client.ProductDelete(ctx, &pb.ProductDeleteRequest{
		Code: code,
	})

	if err != nil {
		log.Printf("func: [%v]:\nerror: [%v]\nexpected: [%v]", funcName, err, expected)
		return
	}

	log.Printf("func: [%v]:\nresponse: [%v]\nexpected: [%v]", funcName, resp, expected)
}
