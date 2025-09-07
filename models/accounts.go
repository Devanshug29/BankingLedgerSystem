package models

import "time"

// CreateAccountRequest represents the payload for creating a new account
type CreateAccountRequest struct {
	AccountNumber string `json:"account_number"`
	Name          string `json:"name" binding:"required"`
	Balance       int64  `json:"balance"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`              // application or HTTP error code
	Message string `json:"message"`           // human-readable error message
	Details string `json:"details,omitempty"` // optional details for debugging
}

type SuccessResponse struct {
	Code int         `json:"code"`           // HTTP status code
	Data interface{} `json:"data,omitempty"` // actual response payload
}

type Account struct {
	ID            int64     `json:"id"`
	AccountNumber string    `json:"account_number"`
	Name          string    `json:"name"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}
