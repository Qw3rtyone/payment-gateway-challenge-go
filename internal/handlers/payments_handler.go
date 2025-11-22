package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/services"
	"github.com/go-chi/chi/v5"
)

type PaymentsHandler struct {
	validator        services.ValidationService
	paymentProcessor services.PaymentService
}

func NewPaymentsHandler(validator services.ValidationService, processor services.PaymentService) *PaymentsHandler {
	return &PaymentsHandler{
		validator:        validator,
		paymentProcessor: processor,
	}
}

// GetHandler returns an http.HandlerFunc that handles HTTP GET requests.
// It retrieves a payment record by its ID from the storage.
// The ID is expected to be part of the URL.
func (h *PaymentsHandler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := chi.URLParam(r, "id")

		response, err := h.paymentProcessor.GetPayment(ctx, id)
		if err != nil {
			if err.Error() != models.ErrPaymentNotFound.Error() {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// PostHandler returns an http.HandlerFunc that handles HTTP POST requests to process payments.
func (h *PaymentsHandler) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		// Parse request body
		var req models.PaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: "Invalid request body",
			})
			return
		}

		// Validate request
		validationErrors := h.validator.ValidatePaymentRequest(ctx, req)
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error:  "Rejected",
				Errors: validationErrors,
			})
			return
		}

		response, err := h.paymentProcessor.CreatePayment(ctx, req)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: "Payment processing failed: " + err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
