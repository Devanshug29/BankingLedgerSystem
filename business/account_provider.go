package service

import (
	"BankingLedgerSystem/models"
	"BankingLedgerSystem/repository"
	"context"
	"errors"
)

// AccountService defines account-related business operations
type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService creates a new service instance with injected repository
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// CreateAccount business logic
func (s *AccountService) CreateAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error) {
	// Business rules can go here
	if req.Balance < 0 {
		return nil, errors.New("balance cannot be negative")
	}

	return s.repo.InsertAccount(ctx, req)
}

// GetAccount business logic
func (s *AccountService) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return s.repo.FindAccountByID(ctx, id)
}
