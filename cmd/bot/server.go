package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"

	apiPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/api"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"google.golang.org/grpc"
)

func TimeoutIntercept(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	deadLineCtx, cancel := context.WithTimeout(ctx, time.Second/2)
	defer cancel()
	waitCh := make(chan struct{}, 1)
	go func() {
		resp, err = handler(deadLineCtx, req)
		waitCh <- struct{}{}
	}()
	select {
	case <-deadLineCtx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "timeout")
	case <-waitCh:
		return resp, err
	}
}

func runGRPCServer(product productPkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(TimeoutIntercept),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminServer(grpcServer, apiPkg.New(product))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
