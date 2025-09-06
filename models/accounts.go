package models

// CreateAccountRequest represents the payload for creating a new account
type CreateAccountRequest struct {
	Name    string  `json:"name" binding:"required"`          // account holder name
	Balance float64 `json:"balance" binding:"required,gte=0"` // initial balance (must be >= 0)
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
