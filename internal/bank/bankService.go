package bank

import "github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"

type Bank interface {
	ProcessPayment(req models.PaymentRequest) (*BankResponse, error)
}
