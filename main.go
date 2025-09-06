package main

import (
	service "BankingLedgerSystem/business"
	"BankingLedgerSystem/repository"
	"context"
	"fmt"

	"BankingLedgerSystem/internal/config"
)

func main() {
	ctx := context.Background()
	config.InitConfigs(ctx, "application.yml")
	cfg := config.GetConfig()
	accountRepo := repository.NewPostgresAccountRepo(db)
	layer := service.NewAccountService(accountRepo)
	fmt.Printf("✅ Loaded config: %+v\n", cfg)
}
