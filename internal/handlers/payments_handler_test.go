package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	mock_services "github.com/cko-recruitment/payment-gateway-challenge-go/internal/services/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockValidator := mock_services.NewMockValidationService(ctrl)
	mockPaymentSvc := mock_services.NewMockPaymentService(ctrl)

	payments := NewPaymentsHandler(mockValidator, mockPaymentSvc)

	r := chi.NewRouter()
	r.Get("/api/payments/{id}", payments.GetHandler())
	r.Post("/api/payments", payments.PostHandler())

	t.Run("GET PaymentFound", func(t *testing.T) {
		payment := &models.PaymentResponse{
			Id:                 "test-id",
			Status:             "Authorized",
			CardNumberLastFour: "5678",
			ExpiryMonth:        10,
			ExpiryYear:         2035,
			Currency:           "GBP",
			Amount:             100,
		}

		req := httptest.NewRequest("GET", "/api/payments/test-id", nil)

		// Mock: payment service returns payment
		mockPaymentSvc.EXPECT().GetPayment(gomock.Any(), "test-id").Return(payment, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"test-id"`)
	})

	t.Run("GET PaymentNotFound", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/api/payments/NonExistingID", nil)
		w := httptest.NewRecorder()

		// Mock: payment not found
		mockPaymentSvc.EXPECT().GetPayment(gomock.Any(), "NonExistingID").Return(nil, models.ErrPaymentNotFound)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("POST CreatePayment Success", func(t *testing.T) {

		createReq := models.PaymentRequest{
			CardNumber:  "1111111111111111",
			ExpiryMonth: 10,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
		}

		createdPayment := &models.PaymentResponse{
			Id:                 "created-id",
			Status:             "Authorized",
			Currency:           "GBP",
			Amount:             100,
			CardNumberLastFour: "1111",
			ExpiryMonth:        10,
			ExpiryYear:         2035,
		}

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/payments", bytes.NewReader(body))

		// Validator OK
		mockValidator.EXPECT().ValidatePaymentRequest(gomock.Any(), createReq).Return(nil)
		// Payment service OK
		mockPaymentSvc.EXPECT().CreatePayment(gomock.Any(), createReq).Return(createdPayment, nil)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"created-id"`)
	})

	t.Run("POST CreatePayment ValidationFails", func(t *testing.T) {

		createReq := models.PaymentRequest{
			CardNumber: "",
		}

		errs := []models.ValidationError{
			{
				Field:   "some field",
				Message: "error message",
			},
		}

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/payments", bytes.NewReader(body))

		// Validator returns error
		mockValidator.EXPECT().ValidatePaymentRequest(gomock.Any(), createReq).Return(errs)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
