package services

import (
	"context"
	"fmt"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/bank"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/utils"
	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, req models.PaymentRequest) (*models.PaymentResponse, error)
	GetPayment(ctx context.Context, id string) (*models.PaymentResponse, error)
}

type paymentService struct {
	storage    repository.PaymentsRepository
	bankClient bank.Bank
}

type Status string

const (
	StatusAuthorized Status = "Authorized"
	StatusDeclined   Status = "Declined"
	StatusRejected   Status = "Rejected"
)

func NewPaymentService(repo repository.PaymentsRepository, bankClient bank.Bank) PaymentService {
	return &paymentService{
		storage:    repo,
		bankClient: bankClient,
	}
}

func (p *paymentService) CreatePayment(ctx context.Context, req models.PaymentRequest) (*models.PaymentResponse, error) {

	// Process payment with bank
	bankResp, err := p.bankClient.ProcessPayment(req)
	if err != nil {
		// If bank returns an error, treat as declined
		return nil, fmt.Errorf("bank processing error: %v", err)
	}

	// Determine payment status based on bank response
	var status Status
	if bankResp.Authorized {
		status = StatusAuthorized
	} else {
		status = StatusDeclined
	}

	// Generate payment ID
	paymentID := uuid.New().String()

	// Get last four digits of card
	lastFour := utils.GetLastFourDigits(req.CardNumber)

	// Create payment record
	payment := models.Payment{
		Id:                 paymentID,
		Status:             string(status),
		CardNumberLastFour: lastFour,
		ExpiryMonth:        req.ExpiryMonth,
		ExpiryYear:         req.ExpiryYear,
		Currency:           req.Currency,
		Amount:             req.Amount,
		AuthorizationCode:  bankResp.AuthorizationCode,
	}

	// Store payment
	paymentErr := p.storage.AddPayment(ctx, payment)
	if paymentErr != nil {
		return nil, fmt.Errorf("failed to store payment: %v", paymentErr)
	}

	return &models.PaymentResponse{
		Id:                 paymentID,
		Status:             string(status),
		CardNumberLastFour: lastFour,
		ExpiryMonth:        req.ExpiryMonth,
		ExpiryYear:         req.ExpiryYear,
		Currency:           req.Currency,
		Amount:             req.Amount,
	}, nil
}

func (p *paymentService) GetPayment(ctx context.Context, id string) (*models.PaymentResponse, error) {

	payment := p.storage.GetPayment(ctx, id)
	if payment == nil {
		return nil, models.ErrPaymentNotFound
	}
	// Convert internal payment to response format
	response := models.PaymentResponse{
		Id:                 payment.Id,
		Status:             payment.Status,
		CardNumberLastFour: payment.CardNumberLastFour,
		ExpiryMonth:        payment.ExpiryMonth,
		ExpiryYear:         payment.ExpiryYear,
		Currency:           payment.Currency,
		Amount:             payment.Amount,
	}
	return &response, nil
}
