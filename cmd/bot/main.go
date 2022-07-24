package main

import (
	botPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/add"
	cmdAddPricePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/add_price"
	cmdDeletePkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/delete"
	cmdHelpPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/list"
	cmdPriceHistoryPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/bot/command/price_history"
	productPkg "gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product"
)

func main() {
	var product productPkg.Interface
	{
		product = productPkg.New()
	}

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
