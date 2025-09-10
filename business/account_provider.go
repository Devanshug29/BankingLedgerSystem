package service

import (
	"BankingLedgerSystem/kafka"
	"BankingLedgerSystem/models"
	"BankingLedgerSystem/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

//go:generate sh -c "sh $(git rev-parse --show-toplevel)/scripts/mock_generator.sh $GOFILE"

// AccountService defines account-related business operations
type AccountService struct {
	repo          repository.AccountRepository
	kafkaProducer *kafka.Producer
}

func NewAccountService(repo repository.AccountRepository, kafkaProducer *kafka.Producer) *AccountService {
	return &AccountService{repo: repo, kafkaProducer: kafkaProducer}
}

func (s *AccountService) CreateAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error) {
	if req.Balance < 0 {
		return nil, errors.New("balance cannot be negative")
	}

	accountNumber := fmt.Sprintf("AC%v", time.Now().Unix())
	req.AccountNumber = accountNumber

	return s.repo.InsertAccount(ctx, req)
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return s.repo.FindAccountByID(ctx, id)
}
func (s *AccountService) PublishTransaction(ctx context.Context, key, payload []byte) error {
	return s.kafkaProducer.Publish(ctx, key, payload)
}
