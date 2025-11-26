package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// TwilioConfig holds Twilio credentials
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// GetTwilioConfig reads Twilio configuration from environment variables
func GetTwilioConfig() *TwilioConfig {
	return &TwilioConfig{
		AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		FromNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
	}
}

// TwilioResponse represents the Twilio API response
type TwilioResponse struct {
	SID         string `json:"sid"`
	Status      string `json:"status"`
	ErrorCode   int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"message,omitempty"`
}

// SendSMS sends an SMS using Twilio API
func SendSMS(to, message string) error {
	config := GetTwilioConfig()

	// Validate configuration
	if config.AccountSID == "" || config.AuthToken == "" || config.FromNumber == "" {
		return fmt.Errorf("twilio credentials not configured")
	}

	// Twilio API URL
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", config.AccountSID)

	// Prepare form data
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", config.FromNumber)
	msgData.Set("Body", message)

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(msgData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.SetBasicAuth(config.AccountSID, config.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var twilioResp TwilioResponse
	if err := json.NewDecoder(resp.Body).Decode(&twilioResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("twilio error (%d): %s", twilioResp.ErrorCode, twilioResp.ErrorMessage)
	}

	fmt.Printf("SMS sent successfully! SID: %s, Status: %s\n", twilioResp.SID, twilioResp.Status)
	return nil
}

// SendOTPSMS sends OTP via SMS
func SendOTPSMS(phoneNumber, otpCode string) error {
	message := fmt.Sprintf("Your OTP verification code is: %s\n\nThis code will expire in 5 minutes.\n\nDo not share this code with anyone.", otpCode)
	return SendSMS(phoneNumber, message)
}
