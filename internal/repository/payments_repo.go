package repository

import (
	"context"
	"sync"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
)

type inMemStore struct {
	mu       sync.RWMutex
	payments map[string]models.Payment
}

func NewPaymentsRepository() PaymentsRepository {
	return &inMemStore{
		payments: make(map[string]models.Payment),
	}
}

func (ps *inMemStore) GetPayment(ctx context.Context, id string) *models.Payment {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if payment, exists := ps.payments[id]; exists {
		return &payment
	}
	return nil
}

func (ps *inMemStore) AddPayment(ctx context.Context, payment models.Payment) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.payments[payment.Id] = payment

	return nil
}
