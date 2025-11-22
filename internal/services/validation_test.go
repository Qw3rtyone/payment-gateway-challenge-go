package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestValidatePaymentRequest_ValidRequest(t *testing.T) {
	req := models.PaymentRequest{
		CardNumber:  "1234567812345678",
		ExpiryMonth: 12,
		ExpiryYear:  time.Now().Year() + 1,
		Currency:    "USD",
		Amount:      1000,
		Cvv:         "123",
	}

	ctx := context.Background()
	errors := NewValidationService().ValidatePaymentRequest(ctx, req)
	assert.Empty(t, errors)
}

func TestValidatePaymentRequest_InvalidCardNumber(t *testing.T) {
	tests := []struct {
		name       string
		cardNumber string
	}{
		{"empty", ""},
		{"too short", "123456789012"},
		{"too long", "12345678901234567890"},
		{"non-numeric", "abcdefghijklmnop"},
	}
	v := NewValidationService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  tt.cardNumber,
				ExpiryMonth: 12,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)
			assert.NotEmpty(t, errors)

			found := false
			for _, err := range errors {
				if err.Field == "card_number" {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected card_number validation error")
		})
	}
}

func TestValidatePaymentRequest_InvalidExpiryMonth(t *testing.T) {
	tests := []struct {
		name  string
		month int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too high", 13},
	}
	v := NewValidationService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: tt.month,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)
			assert.NotEmpty(t, errors)

			found := false
			for _, err := range errors {
				if err.Field == "expiry_month" {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected expiry_month validation error")
		})
	}
}

func TestValidatePaymentRequest_ExpiryDateInPast(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		month int
		year  int
	}{
		{"last year", 12, now.Year() - 1},
		{"last month", int(now.Month()) - 1, now.Year()},
	}
	v := NewValidationService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.month < 1 {
				t.Skip("Skipping test case with invalid month")
			}

			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: tt.month,
				ExpiryYear:  tt.year,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         "123",
			}
			fmt.Printf("req: %v\n", req)
			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)
			assert.NotEmpty(t, errors)
		})
	}
}

func TestValidatePaymentRequest_InvalidCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
	}{
		{"empty", ""},
		{"too short", "US"},
		{"too long", "USDD"},
		{"unsupported", "JPY"},
	}
	v := NewValidationService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: 12,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    tt.currency,
				Amount:      1000,
				Cvv:         "123",
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)
			assert.NotEmpty(t, errors)

			found := false
			for _, err := range errors {
				if err.Field == "currency" {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected currency validation error")
		})
	}
}

func TestValidatePaymentRequest_SupportedCurrencies(t *testing.T) {
	currencies := []string{"USD", "GBP", "EUR"}
	v := NewValidationService()
	for _, currency := range currencies {
		t.Run(currency, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: 12,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    currency,
				Amount:      1000,
				Cvv:         "123",
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)

			for _, err := range errors {
				assert.NotEqual(t, "currency", err.Field, "Currency %s should be valid", currency)
			}
		})
	}
}

func TestValidatePaymentRequest_InvalidAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount int
	}{
		{"zero", 0},
		{"negative", -100},
	}
	v := NewValidationService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: 12,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    "USD",
				Amount:      tt.amount,
				Cvv:         "123",
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)
			assert.NotEmpty(t, errors)

			found := false
			for _, err := range errors {
				if err.Field == "amount" {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected amount validation error")
		})
	}
}

func TestValidatePaymentRequest_InvalidCvv(t *testing.T) {
	tests := []struct {
		name string
		cvv  string
	}{
		{"empty", ""},
		{"too short", "12"},
		{"too long", "12345"},
		{"non-numeric", "12a"},
	}
	v := NewValidationService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: 12,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         tt.cvv,
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)
			assert.NotEmpty(t, errors)

			found := false
			for _, err := range errors {
				if err.Field == "cvv" {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected cvv validation error")
		})
	}
}

func TestValidatePaymentRequest_ValidCvv(t *testing.T) {
	tests := []string{"123", "1234"}
	v := NewValidationService()
	for _, cvv := range tests {
		t.Run(cvv, func(t *testing.T) {
			req := models.PaymentRequest{
				CardNumber:  "1234567812345678",
				ExpiryMonth: 12,
				ExpiryYear:  time.Now().Year() + 1,
				Currency:    "USD",
				Amount:      1000,
				Cvv:         cvv,
			}

			ctx := context.Background()
			errors := v.ValidatePaymentRequest(ctx, req)

			for _, err := range errors {
				assert.NotEqual(t, "cvv", err.Field, "CVV %s should be valid", cvv)
			}
		})
	}
}
