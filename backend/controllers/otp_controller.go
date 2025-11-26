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
	var smsStatus string = "not_sent"
	if req.Phone != "" {
		// Check if Twilio is configured
		if os.Getenv("TWILIO_ACCOUNT_SID") != "" && os.Getenv("TWILIO_AUTH_TOKEN") != "" && os.Getenv("TWILIO_PHONE_NUMBER") != "" {
			fmt.Printf("\n📱 SMS Sending Process Started...\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("📤 Destination: %s\n", req.Phone)
			fmt.Printf("🔐 OTP Code: %s\n", otpCode)
			fmt.Printf("🔑 Twilio SID: %s...\n", os.Getenv("TWILIO_ACCOUNT_SID")[:10])
			fmt.Printf("📞 From Number: %s\n", os.Getenv("TWILIO_PHONE_NUMBER"))
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			// Send via Twilio
			if err := utils.SendOTPSMS(req.Phone, otpCode); err != nil {
				fmt.Printf("\n❌ SMS Delivery Failed!\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("\n💡 Possible Reasons:\n")
				fmt.Printf("   1. Phone number not verified (Trial Account)\n")
				fmt.Printf("   2. Invalid Twilio credentials\n")
				fmt.Printf("   3. Insufficient Twilio credits\n")
				fmt.Printf("   4. Wrong phone number format\n")
				fmt.Printf("\n🔧 Solutions:\n")
				fmt.Printf("   1. Verify phone at: https://console.twilio.com/\n")
				fmt.Printf("   2. Check .env Twilio credentials\n")
				fmt.Printf("   3. Ensure phone format: +919876543210\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
				smsStatus = "failed"
			} else {
				fmt.Printf("\n✅ SMS Sent Successfully!\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
				fmt.Printf("✓ Message queued for delivery\n")
				fmt.Printf("✓ User will receive SMS shortly\n")
				fmt.Printf("✓ Check Twilio Console for delivery status\n")
				fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
				smsStatus = "sent"
			}
		} else {
			fmt.Printf("\n⚠️  Twilio Configuration Missing!\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Required in .env file:\n")
			fmt.Printf("  TWILIO_ACCOUNT_SID=ACxxxxxxxxxx\n")
			fmt.Printf("  TWILIO_AUTH_TOKEN=your_token\n")
			fmt.Printf("  TWILIO_PHONE_NUMBER=+1234567890\n")
			fmt.Printf("\n📚 Setup Guide: TWILIO_SETUP.md\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
			smsStatus = "twilio_not_configured"
		}
	}

	// Send OTP via Email if email is provided
	if req.Email != "" {
		fmt.Printf("\n📧 Email OTP Feature\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Status: Not yet implemented\n")
		fmt.Printf("TODO: Integrate SendGrid/AWS SES/Gmail\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")
	}

	// Log for development
	fmt.Printf("\n═══════════════════════════════════════════\n")
	fmt.Printf("         🔐 OTP GENERATED                 \n")
	fmt.Printf("═══════════════════════════════════════════\n")
	fmt.Printf("OTP ID:      %s\n", otp.ID)
	fmt.Printf("OTP Code:    %s\n", otpCode)
	fmt.Printf("Email:       %s\n", req.Email)
	fmt.Printf("Phone:       %s\n", req.Phone)
	fmt.Printf("SMS Status:  %s\n", smsStatus)
	fmt.Printf("Expires At:  %s\n", otp.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Valid For:   5 minutes\n")
	fmt.Printf("═══════════════════════════════════════════\n\n")

	// Response data
	responseData := gin.H{
		"otp_id":     otp.ID,
		"expires_at": otp.ExpiresAt,
		"sms_status": smsStatus,
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

	fmt.Printf("\n🔍 OTP Verification Attempt\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("OTP ID: %s\n", req.OTPID)
	fmt.Printf("Code Provided: %s\n", req.OTPCode)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Find OTP record
	var otp models.OTP
	if err := config.DB.Where("id = ?", req.OTPID).First(&otp).Error; err != nil {
		fmt.Printf("❌ OTP not found in database\n\n")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "OTP not found",
		})
		return
	}

	// Check if OTP is already verified
	if otp.IsVerified {
		fmt.Printf("❌ OTP already used\n\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP already verified",
		})
		return
	}

	// Check if OTP has expired
	if time.Now().After(otp.ExpiresAt) {
		fmt.Printf("❌ OTP expired at: %s\n\n", otp.ExpiresAt.Format("2006-01-02 15:04:05"))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "OTP has expired",
		})
		return
	}

	// Check maximum attempts (3 attempts allowed)
	if otp.AttemptCount >= 3 {
		fmt.Printf("❌ Maximum attempts exceeded\n\n")
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
		fmt.Printf("❌ Invalid OTP code. Attempts remaining: %d\n\n", 3-otp.AttemptCount)
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

	fmt.Printf("\n✅ OTP VERIFIED SUCCESSFULLY!\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Phone: %s\n", user.Phone)
	fmt.Printf("Verified At: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

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

	fmt.Printf("\n🔄 OTP Resend Request\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("OTP ID: %s\n", req.OTPID)
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Find old OTP record
	var oldOTP models.OTP
	if err := config.DB.Where("id = ?", req.OTPID).First(&oldOTP).Error; err != nil {
		fmt.Printf("❌ OTP not found\n\n")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "OTP not found",
		})
		return
	}

	// Check if already verified
	if oldOTP.IsVerified {
		fmt.Printf("❌ OTP already verified\n\n")
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
	var smsStatus string = "not_sent"
	if oldOTP.Phone != "" {
		if os.Getenv("TWILIO_ACCOUNT_SID") != "" {
			fmt.Printf("📤 Resending SMS to: %s\n", oldOTP.Phone)
			if err := utils.SendOTPSMS(oldOTP.Phone, otpCode); err != nil {
				fmt.Printf("❌ Failed to resend SMS: %v\n\n", err)
				smsStatus = "failed"
			} else {
				fmt.Printf("✅ SMS resent successfully\n\n")
				smsStatus = "sent"
			}
		}
	}

	// Log for development
	fmt.Printf("\n═══════════════════════════════════════════\n")
	fmt.Printf("         🔁 OTP RESENT                    \n")
	fmt.Printf("═══════════════════════════════════════════\n")
	fmt.Printf("New OTP ID:  %s\n", newOTP.ID)
	fmt.Printf("OTP Code:    %s\n", otpCode)
	fmt.Printf("Phone:       %s\n", oldOTP.Phone)
	fmt.Printf("SMS Status:  %s\n", smsStatus)
	fmt.Printf("Expires At:  %s\n", newOTP.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("═══════════════════════════════════════════\n\n")

	// Response data
	responseData := gin.H{
		"otp_id":     newOTP.ID,
		"expires_at": newOTP.ExpiresAt,
		"sms_status": smsStatus,
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
