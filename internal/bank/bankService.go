package bank

import (
	"context"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
)

type Bank interface {
	ProcessPayment(ctx context.Context, req models.PaymentRequest) (*BankResponse, error)
}
