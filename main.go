package main

import (
	"BankingLedgerSystem/api"
	controller "BankingLedgerSystem/api/v1"
	service "BankingLedgerSystem/business"
	"BankingLedgerSystem/internal/flags"
	"BankingLedgerSystem/kafka"
	"BankingLedgerSystem/repository"
	"context"
	"fmt"
	"log"

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
		log.Panicf("failed to initialize DB: %v", err)
	}

	monogCOnnectionProvider := repository.NewMongoProvider(cfg)
	err = monogCOnnectionProvider.InitMongo(ctx)
	if err != nil {
		log.Panicf("failed to initialize Mongo DB: %v", err)
	}

	consumer := kafka.NewConsumer(cfg, connectionProvider)
	defer consumer.Close()

	// Start consumer in background
	go func() {
		log.Println("Kafka consumer started...")
		consumer.Start(ctx)
	}()
	producer := kafka.NewProducer(cfg)
	accountRepo := repository.NewPostgresAccountRepo(connectionProvider)
	layer := service.NewAccountService(accountRepo, producer)
	accountController := controller.NewAccountController(layer)

	router := api.GetRouter(ctx)
	api.RegisterRoutes(router, cfg, accountController)

	appServer := &server.Server{
		Addr:    fmt.Sprintf(":%d", flags.GetDeploymentMode().Port()),
		Handler: router,
	}
	if err := appServer.ListenAndServe(); err != nil {
		log.Panicf("failed to initialize Router: %v", err)
	}

}
