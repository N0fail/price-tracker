package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
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
	listener, err := net.Listen("tcp", config.GrpcPort)
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

func runREST() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", gmux)

	// Register Swagger Handler
	fs := http.FileServer(http.Dir("./swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, gmux, config.GrpcPort, opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(config.RESTPort, mux); err != nil {
		panic(err)
	}
}

func headerMatcherREST(key string) (string, bool) {
	switch key {
	case "Custom":
		return key, true
	default:
		return key, false
	}
}

func runSwagger() {

}
