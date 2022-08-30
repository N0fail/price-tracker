package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/kafka"
	botPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/add"
	cmdAddPricePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/add_price"
	cmdDeletePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/delete"
	cmdHelpPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/list"
	cmdPriceHistoryPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/price_history"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	useCache := flag.Bool("cache", false, "run using cache")
	flag.Parse()

	var pool *pgxpool.Pool
	if !*useCache {
		logrus.Info("Run using database")
		// connection string
		psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DbHost, config.DbPort, config.DbUser, config.DbPassword, config.DbName)

		// connect to database
		var err error
		pool, err = pgxpool.Connect(ctx, psqlConn)
		if err != nil {
			log.Fatal("can't connect to database", err)
		}
		defer pool.Close()

		if err := pool.Ping(ctx); err != nil {
			log.Fatal("ping database error", err)
		}

		// настраиваем
		poolConfig := pool.Config()
		poolConfig.MaxConnIdleTime = config.DbMaxConnIdleTime
		poolConfig.MaxConnLifetime = config.DbMaxConnLifetime
		poolConfig.MinConns = config.DbMinConns
		poolConfig.MaxConns = config.DbMaxConns
	} else {
		logrus.Info("Run using cache")
	}

	var product productPkg.Interface
	{
		product = productPkg.New(pool)
	}
	//go runBot(product)
	go runREST()
	go kafka.Run(product)
	go http.ListenAndServe("127.0.0.1:8200", nil)
	runGRPCServer(product)
}

func runBot(product productPkg.Interface) {
	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()
		bot.RegisterCommand(cmdAddPkg.New(product))
		bot.RegisterCommand(cmdAddPricePkg.New(product))
		bot.RegisterCommand(cmdDeletePkg.New(product))
		bot.RegisterCommand(cmdListPkg.New(product))
		bot.RegisterCommand(cmdPriceHistoryPkg.New(product))
		bot.RegisterCommand(cmdHelpPkg.New(bot.GetCommands()))
	}

	bot.Run()
}
