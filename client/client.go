package main

import (
	"context"
	"log"

	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conns, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAdminClient(conns)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	list, err := client.ProductList(ctx, &pb.ProductListRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response: [%v]", list)

	create, err := client.ProductCreate(ctx, &pb.ProductCreateRequest{
		Code: "1",
		Name: "333",
	})
	if err != nil {
		log.Print(err)
	}

	log.Printf("response: [%v]", create)

	create, err = client.ProductCreate(ctx, &pb.ProductCreateRequest{
		Code: "1",
		Name: "4444",
	})
	if err != nil {
		log.Print(err)
	}
	log.Printf("response: [%v]", create)

	create, err = client.ProductCreate(ctx, &pb.ProductCreateRequest{
		Code: "2",
		Name: "5555",
	})
	if err != nil {
		log.Print(err)
	}
	log.Printf("response: [%v]", create)

	list, err = client.ProductList(ctx, &pb.ProductListRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("response: [%v]", list)

}
