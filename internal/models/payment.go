package models

import "errors"

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

type PaymentRequest struct {
	CardNumber  string `json:"card_number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	Currency    string `json:"currency"`
	Amount      int    `json:"amount"`
	Cvv         string `json:"cvv"`
}

type PaymentResponse struct {
	Id                 string `json:"id"`
	Status             string `json:"status"`
	CardNumberLastFour string `json:"card_number_last_four"`
	ExpiryMonth        int    `json:"expiry_month"`
	ExpiryYear         int    `json:"expiry_year"`
	Currency           string `json:"currency"`
	Amount             int    `json:"amount"`
}

// Payment represents the internal storage model
type Payment struct {
	Id                 string
	Status             string
	CardNumberLastFour string
	ExpiryMonth        int
	ExpiryYear         int
	Currency           string
	Amount             int
	AuthorizationCode  string
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error  string            `json:"error"`
	Errors []ValidationError `json:"errors,omitempty"`
}
