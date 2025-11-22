package services

import (
	"context"
	"regexp"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
)

var (
	// Supported currencies
	supportedCurrencies = map[string]bool{
		"USD": true,
		"GBP": true,
		"EUR": true,
	}

	numericRegex = regexp.MustCompile(`^[0-9]+$`)
)

type ValidationService interface {
	ValidatePaymentRequest(ctx context.Context, req models.PaymentRequest) []models.ValidationError
}

type validationService struct{}

func NewValidationService() ValidationService {
	return &validationService{}
}

// ValidatePaymentRequest validates all fields in a payment request
func (v *validationService) ValidatePaymentRequest(ctx context.Context, req models.PaymentRequest) []models.ValidationError {
	return concatErrors(
		validateCardNumber(req.CardNumber),
		validateExpiryDate(req.ExpiryMonth, req.ExpiryYear),
		validateAmount(req.Amount),
		validateCurrency(req.Currency),
		validateCvv(req.Cvv),
	)
}

func validateAmount(amount int) []models.ValidationError {
	var errors []models.ValidationError

	if amount <= 0 {
		errors = append(errors, models.ValidationError{
			Field:   "amount",
			Message: "amount must be a positive integer",
		})
	}

	return errors
}

func validateCardNumber(cardNumber string) []models.ValidationError {
	var errors []models.ValidationError

	if cardNumber == "" {
		errors = append(errors, models.ValidationError{
			Field:   "card_number",
			Message: "card number is required",
		})
	} else {
		if len(cardNumber) < 14 || len(cardNumber) > 19 {
			errors = append(errors, models.ValidationError{
				Field:   "card_number",
				Message: "card number must be between 14-19 characters long",
			})
		}
		if !numericRegex.MatchString(cardNumber) {
			errors = append(errors, models.ValidationError{
				Field:   "card_number",
				Message: "card number must contain only numeric characters",
			})
		}
	}
	return errors
}

func validateExpiryDate(month int, year int) []models.ValidationError {
	var errors []models.ValidationError
	errors = append(errors, validateExpiryYear(year)...)
	errors = append(errors, validateExpiryMonth(month)...)

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	if year == currentYear && month < currentMonth {
		errors = append(errors, models.ValidationError{
			Field:   "expiry_month",
			Message: "expiry date must be in the future",
		})
	}
	return errors
}
func validateExpiryMonth(month int) []models.ValidationError {
	var errors []models.ValidationError
	if month < 1 || month > 12 {
		errors = append(errors, models.ValidationError{
			Field:   "expiry_month",
			Message: "expiry month must be between 1-12",
		})
	}

	return errors
}

func validateExpiryYear(year int) []models.ValidationError {
	var errors []models.ValidationError

	if year == 0 {
		errors = append(errors, models.ValidationError{
			Field:   "expiry_year",
			Message: "expiry year is required",
		})
	} else {
		now := time.Now()
		currentYear := now.Year()

		// Check if the expiry date is in the future
		if year < currentYear {
			errors = append(errors, models.ValidationError{
				Field:   "expiry_year",
				Message: "expiry year must be in the future",
			})
		}
	}

	return errors
}

func validateCurrency(currency string) []models.ValidationError {
	var errors []models.ValidationError
	if currency == "" {
		errors = append(errors, models.ValidationError{
			Field:   "currency",
			Message: "currency is required",
		})
	} else {
		if len(currency) != 3 {
			errors = append(errors, models.ValidationError{
				Field:   "currency",
				Message: "currency must be 3 characters",
			})
		}
		if !supportedCurrencies[currency] {
			errors = append(errors, models.ValidationError{
				Field:   "currency",
				Message: "currency must be one of: USD, GBP, EUR",
			})
		}
	}
	return errors
}

func validateCvv(cvv string) []models.ValidationError {
	var errors []models.ValidationError
	if cvv == "" {
		errors = append(errors, models.ValidationError{
			Field:   "cvv",
			Message: "cvv is required",
		})
	} else {
		if len(cvv) < 3 || len(cvv) > 4 {
			errors = append(errors, models.ValidationError{
				Field:   "cvv",
				Message: "cvv must be 3-4 characters long",
			})
		}
		if !numericRegex.MatchString(cvv) {
			errors = append(errors, models.ValidationError{
				Field:   "cvv",
				Message: "cvv must contain only numeric characters",
			})
		}
	}
	return errors
}

func concatErrors(slicesOfErrs ...[]models.ValidationError) []models.ValidationError {
	var result []models.ValidationError
	for _, s := range slicesOfErrs {
		if len(s) > 0 {
			result = append(result, s...)
		}
	}
	return result
}
