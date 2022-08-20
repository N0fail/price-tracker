//go:build integration
// +build integration

package tests

import (
	"gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"gitlab.ozon.dev/N0fail/price-tracker/tests/config"
	"gitlab.ozon.dev/N0fail/price-tracker/tests/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	AdminClient api.AdminClient
	Db          *postgres.TDB
)

func init() {
	cfg, err := config.FromEnv()

	conn, err := grpc.Dial(cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(3*time.Second))
	if err != nil {
		panic(err)
	}
	AdminClient = api.NewAdminClient(conn)

	Db = postgres.NewFromEnv()
}
