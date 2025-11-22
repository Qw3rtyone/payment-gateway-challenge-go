package utils

import (
	"fmt"
	"strconv"
)

// GetLastFourDigits extracts the last 4 digits from a card number
func GetLastFourDigits(cardNumber string) string {
	if len(cardNumber) >= 4 {
		return cardNumber[len(cardNumber)-4:]
	}
	return cardNumber
}

// FormatExpiryDate formats expiry month and year as MM/YY for bank API
func FormatExpiryDate(month, year int) string {
	// Convert year to 2-digit format
	yearStr := strconv.Itoa(year)
	if len(yearStr) == 4 {
		yearStr = yearStr[2:]
	}
	return fmt.Sprintf("%02d/%s", month, yearStr)
}
