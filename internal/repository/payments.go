package repository

import (
	"context"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
)

type PaymentsRepository interface {
	GetPayment(ctx context.Context, id string) *models.Payment
	AddPayment(ctx context.Context, payment models.Payment) error
}
