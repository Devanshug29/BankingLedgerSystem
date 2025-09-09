package models

import "time"

// CreateAccountRequest represents the payload for creating a new account
type CreateAccountRequest struct {
	AccountNumber string `json:"account_number"`
	Name          string `json:"name" binding:"required"`
	Balance       int64  `json:"balance"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type Account struct {
	ID            int64     `json:"id"`
	AccountNumber string    `json:"account_number"`
	Name          string    `json:"name"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionRequest struct {
	AccountNumber string `json:"accountNumber" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=deposit withdraw"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}
