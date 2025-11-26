package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"otp-backend/config"
	"otp-backend/models"
	"otp-backend/utils"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GenerateOTPRequest represents the request body for OTP generation
type GenerateOTPRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
	Phone string `json:"phone" binding:"omitempty,min=10,max=15"`
}

// VerifyOTPRequest represents the request body for OTP verification
type VerifyOTPRequest struct {
	OTPID   string `json:"otp_id" binding:"required"`
	OTPCode string `json:"otp_code" binding:"required,len=6"`
}

// ResendOTPRequest represents the request body for resending OTP
type ResendOTPRequest struct {
	OTPID string `json:"otp_id" binding:"required"`
}

// GenerateOTP generates a new OTP and sends it to the user
func GenerateOTP(c *gin.Context) {
	var req GenerateOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validate that at least email or phone is provided
	if req.Email == "" && req.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Either email or phone number is required",
		})
		return
	}

	// Check rate limiting (max 3 OTP requests per hour)
	var count int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	query := config.DB.Model(&models.OTP{}).Where("created_at > ?", oneHourAgo)

	if req.Email != "" {
		query = query.Where("email = ?", req.Email)
	} else {
		query = query.Where("phone = ?", req.Phone)
	}

	query.Count(&count)

	if count >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success": false,
			"message": "Too many OTP requests. Please try again after an hour",
		})
		return
	}

	// Generate 6-digit OTP
	otpCode, err := generateSecureOTP(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate OTP",
			"error":   err.Error(),
		})
		return
	}

	// Create OTP record
	otp := models.OTP{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Phone:        req.Phone,
		OTPCode:      otpCode,
		IsVerified:   false,
		AttemptCount: 0,
		ExpiresAt:    time.Now().Add(5 * time.Minute), // OTP valid for 5 minutes
	}

	if err := config.DB.Create(&otp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save OTP",
			"error":   err.Error(),
		})
		return
	}

	// Send OTP via SMS if phone number is provided
	if req.Phone != "" {
		// Check if Twilio is configured
		if os.Getenv("TWILIO_ACCOUNT_SID") != "" {
			// Send via Twilio
			if err := utils.SendOTPSMS(req.Phone, otpCode); err != nil {
				fmt.Printf("Failed to send SMS via Twilio: %v\n", err)
				// Don't fail the request, just log the error
			}
		} else {
			fmt.Printf("\n⚠️  Twilio not configured. Set TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_PHONE_NUMBER\n")
		}
	}

	// Send OTP via Email if email is provided
	if req.Email != "" {
		// TODO: Implement email sending (SendGrid, AWS SES, etc.)
		fmt.Printf("\n📧 Email OTP sending not yet implemented\n")
	}

	// Log for development
	fmt.Printf("\n=== OTP Generated ===\n")
	fmt.Printf("OTP ID: %s\n", otp.ID)
	fmt.Printf("OTP Code: %s\n", otpCode)
	fmt.Printf("Email: %s\n", req.Email)
	fmt.Printf("Phone: %s\n", req.Phone)
	fmt.Printf("Expires At: %s\n", otp.ExpiresAt.Format(time.RFC3339))
	fmt.Printf("===================\n\n")

	// Response data
	responseData := gin.H{
		"otp_id":     otp.ID,
		"expires_at": otp.ExpiresAt,
	}

	// Only include OTP code in development mode
	if os.Getenv("ENVIRONMENT") != "production" {
		responseData["otp_code"] = otpCode
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP sent successfully",
		"data":    responseData,
	})
}

// VerifyOTP verifies the provided OTP code
func VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Find OTP record
	var otp models.OTP
	if err := config.DB.Where("id = ?", req.OTPID).First(&otp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "OTP not found",
		})
		return
	}

	// Check if OTP is already verified
	if otp.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP already verified",
		})
		return
	}

	// Check if OTP has expired
	if time.Now().After(otp.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP has expired",
		})
		return
	}

	// Check maximum attempts (3 attempts allowed)
	if otp.AttemptCount >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Maximum verification attempts exceeded",
		})
		return
	}

	// Increment attempt count
	otp.AttemptCount++
	config.DB.Save(&otp)

	// Verify OTP code
	if otp.OTPCode != req.OTPCode {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Invalid OTP code. %d attempts remaining", 3-otp.AttemptCount),
		})
		return
	}

	// Mark OTP as verified
	now := time.Now()
	otp.IsVerified = true
	otp.VerifiedAt = &now
	config.DB.Save(&otp)

	// Update or create user record
	var user models.User
	var userExists bool

	if otp.Email != "" {
		userExists = config.DB.Where("email = ?", otp.Email).First(&user).Error == nil
	} else if otp.Phone != "" {
		userExists = config.DB.Where("phone = ?", otp.Phone).First(&user).Error == nil
	}

	if userExists {
		// Update existing user
		if otp.Email != "" {
			user.IsEmailVerified = true
		}
		if otp.Phone != "" {
			user.IsPhoneVerified = true
		}
		config.DB.Save(&user)
	} else {
		// Create new user
		user = models.User{
			ID:              uuid.New().String(),
			Email:           otp.Email,
			Phone:           otp.Phone,
			IsEmailVerified: otp.Email != "",
			IsPhoneVerified: otp.Phone != "",
		}
		config.DB.Create(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP verified successfully",
		"data": gin.H{
			"verified":  true,
			"user_id":   user.ID,
			"email":     user.Email,
			"phone":     user.Phone,
			"timestamp": now,
		},
	})
}

// ResendOTP resends an OTP
func ResendOTP(c *gin.Context) {
	var req ResendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Find old OTP record
	var oldOTP models.OTP
	if err := config.DB.Where("id = ?", req.OTPID).First(&oldOTP).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "OTP not found",
		})
		return
	}

	// Check if already verified
	if oldOTP.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP already verified",
		})
		return
	}

	// Generate new OTP
	otpCode, err := generateSecureOTP(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate OTP",
			"error":   err.Error(),
		})
		return
	}

	// Create new OTP record
	newOTP := models.OTP{
		ID:           uuid.New().String(),
		Email:        oldOTP.Email,
		Phone:        oldOTP.Phone,
		OTPCode:      otpCode,
		IsVerified:   false,
		AttemptCount: 0,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}

	if err := config.DB.Create(&newOTP).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save OTP",
			"error":   err.Error(),
		})
		return
	}

	// Resend OTP via SMS if phone number is provided
	if oldOTP.Phone != "" {
		if os.Getenv("TWILIO_ACCOUNT_SID") != "" {
			if err := utils.SendOTPSMS(oldOTP.Phone, otpCode); err != nil {
				fmt.Printf("Failed to send SMS via Twilio: %v\n", err)
			}
		}
	}

	// Log for development
	fmt.Printf("\n=== OTP Resent ===\n")
	fmt.Printf("New OTP ID: %s\n", newOTP.ID)
	fmt.Printf("OTP Code: %s\n", otpCode)
	fmt.Printf("Expires At: %s\n", newOTP.ExpiresAt.Format(time.RFC3339))
	fmt.Printf("==================\n\n")

	// Response data
	responseData := gin.H{
		"otp_id":     newOTP.ID,
		"expires_at": newOTP.ExpiresAt,
	}

	// Only include OTP code in development mode
	if os.Getenv("ENVIRONMENT") != "production" {
		responseData["otp_code"] = otpCode
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP resent successfully",
		"data":    responseData,
	})
}

// generateSecureOTP generates a cryptographically secure random OTP
func generateSecureOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
