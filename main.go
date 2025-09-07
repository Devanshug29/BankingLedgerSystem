package main

import (
	"BankingLedgerSystem/api"
	controller "BankingLedgerSystem/api/v1"
	service "BankingLedgerSystem/business"
	"BankingLedgerSystem/internal/flags"
	"BankingLedgerSystem/repository"
	"context"
	"fmt"

	_ "BankingLedgerSystem/docs"
	"BankingLedgerSystem/internal/config"
	server "net/http"
)

func main() {
	ctx := context.Background()
	config.InitConfigs(ctx, "application.yml")
	cfg := config.GetConfig()

	connectionProvider := repository.NewConnectionProvider(cfg)
	err := connectionProvider.InitDB(ctx, "ledgerdb")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize DB: %v", err))
	}

	accountRepo := repository.NewPostgresAccountRepo(connectionProvider)
	layer := service.NewAccountService(accountRepo)
	accountController := controller.NewAccountController(layer)

	router := api.GetRouter(ctx)
	api.RegisterRoutes(router, cfg, accountController)

	appServer := &server.Server{
		Addr:    fmt.Sprintf(":%d", flags.GetDeploymentMode().Port()),
		Handler: router,
	}
	if err := appServer.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("failed to initialize Router: %v", err))
	}

	fmt.Printf("Loaded config: %+v\n", cfg)
}
