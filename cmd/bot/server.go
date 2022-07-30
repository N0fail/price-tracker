package main

import (
	"log"
	"net"

	apiPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/api"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"google.golang.org/grpc"
)

func runGRPCServer(product productPkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(product))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
