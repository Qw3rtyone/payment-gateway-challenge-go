package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type pong struct {
	Message string `json:"message"`
}

// PingHandler returns an http.HandlerFunc that handles HTTP Ping GET requests.
func (a *Api) PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(pong{Message: "pong"}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// SwaggerHandler returns an http.HandlerFunc that handles HTTP Swagger related requests.
func (a *Api) SwaggerHandler() http.HandlerFunc {
	return httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", docs.SwaggerInfo.Host)),
	)
}

// PostPaymentHandler returns an http.HandlerFunc that handles Payments POST requests.
//
//	@Summary		Process a payment
//	@Description	Processes a card payment through the payment gateway
//	@Tags			payments
//	@Accept			json
//	@Produce		json
//	@Param			payment	body		models.PaymentRequest	true	"Payment Request"
//	@Success		200		{object}	models.PaymentResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/api/payments [post]
func (a *Api) PostPaymentHandler() http.HandlerFunc {
	return a.paymentsHandlers.PostHandler()
}

// GetPaymentHandler returns an http.HandlerFunc that handles Payments GET requests.
//
//	@Summary		Retrieve payment details
//	@Description	Retrieves details of a previously made payment by its ID
//	@Tags			payments
//	@Produce		json
//	@Param			id	path		string	true	"Payment ID"
//	@Success		200	{object}	models.PaymentResponse
//	@Failure		404
//	@Router			/api/payments/{id} [get]
func (a *Api) GetPaymentHandler() http.HandlerFunc {
	return a.paymentsHandlers.GetHandler()
}
