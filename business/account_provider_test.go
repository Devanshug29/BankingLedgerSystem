package service

import (
	"BankingLedgerSystem/internal/config"
	"BankingLedgerSystem/mocks/kafka_mocks"
	"BankingLedgerSystem/mocks/repository_mocks"
	"BankingLedgerSystem/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_Success(t *testing.T) {
	ctx := context.Background()
	config.InitConfigs(ctx, "application.yml")

	mockRepo := new(repository_mocks.MockAccountRepository)
	mockProducer := new(kafka_mocks.MockProducerInterface)

	svc := NewAccountService(mockRepo, mockProducer)

	req := models.CreateAccountRequest{Balance: 100}
	expected := &models.Account{
		AccountNumber: "AC123",
		Balance:       100,
	}

	acc, err := svc.CreateAccount(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, expected.Balance, acc.Balance)
	assert.Contains(t, acc.AccountNumber, "AC")

}

func TestCreateAccount_NegativeBalance(t *testing.T) {
	mockRepo := new(repository_mocks.MockAccountRepository)
	mockProducer := new(kafka_mocks.MockProducerInterface)
	svc := NewAccountService(mockRepo, mockProducer)
	req := models.CreateAccountRequest{Balance: -50}
	ctx := context.Background()

	acc, err := svc.CreateAccount(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, acc)
	assert.Equal(t, "balance cannot be negative", err.Error())
}

func TestGetAccount_Success(t *testing.T) {
	mockRepo := new(repository_mocks.MockAccountRepository)
	mockProducer := new(kafka_mocks.MockProducerInterface)
	svc := NewAccountService(mockRepo, mockProducer)

	expected := &models.Account{AccountNumber: "AC999", Balance: 500}

	ctx := context.Background()
	acc, err := svc.GetAccount(ctx, "AC999")

	assert.NoError(t, err)
	assert.Equal(t, expected, acc)
}

func TestGetAccount_NotFound(t *testing.T) {
	mockRepo := new(repository_mocks.MockAccountRepository)
	mockProducer := new(kafka_mocks.MockProducerInterface)
	svc := NewAccountService(mockRepo, mockProducer)

	ctx := context.Background()
	acc, err := svc.GetAccount(ctx, "AC404")

	assert.Error(t, err)
	assert.Nil(t, acc)
}

func TestPublishTransaction(t *testing.T) {
	mockRepo := new(repository_mocks.MockAccountRepository)
	mockProducer := new(kafka_mocks.MockProducerInterface)
	svc := NewAccountService(mockRepo, mockProducer)

	key := []byte("key1")
	payload := []byte(`{"type":"deposit","amount":100}`)

	ctx := context.Background()
	err := svc.PublishTransaction(ctx, key, payload)

	assert.NoError(t, err)
}
