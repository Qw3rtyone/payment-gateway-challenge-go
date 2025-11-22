package bank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/utils"
)

const (
	defaultBankURL = "http://localhost:8080"
)

// BankRequest represents the request format expected by the bank simulator
type BankRequest struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"` // Format: MM/YY
	Currency   string `json:"currency"`
	Amount     int    `json:"amount"`
	Cvv        string `json:"cvv"`
}

// BankResponse represents the response from the bank simulator
type BankResponse struct {
	Authorized        bool   `json:"authorized"`
	AuthorizationCode string `json:"authorization_code"`
}

// Client handles communication with the acquiring bank
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new bank client
func NewClient(url *string) *Client {
	baseUrl := defaultBankURL
	if url != nil && *url != "" {
		baseUrl = *url
	}

	return &Client{
		baseURL: baseUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ProcessPayment sends a payment request to the acquiring bank
func (c *Client) ProcessPayment(ctx context.Context, req models.PaymentRequest) (*BankResponse, error) {
	// Convert payment request to bank request format
	bankReq := BankRequest{
		CardNumber: req.CardNumber,
		ExpiryDate: utils.FormatExpiryDate(req.ExpiryMonth, req.ExpiryYear),
		Currency:   req.Currency,
		Amount:     req.Amount,
		Cvv:        req.Cvv,
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(bankReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	u, err := url.JoinPath(c.baseURL, "payments")
	if err != nil {
		return nil, fmt.Errorf("failed to create bank URL: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", u, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to bank: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusServiceUnavailable {
			return nil, fmt.Errorf("bank service unavailable")
		}
		return nil, fmt.Errorf("bank returned error status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var bankResp BankResponse
	if err := json.Unmarshal(body, &bankResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &bankResp, nil
}
