package routes

import (
	"otp-backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterOTPRoutes registers all OTP-related routes
func RegisterOTPRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		otp := api.Group("/otp")
		{
			otp.POST("/generate", controllers.GenerateOTP)
			otp.POST("/verify", controllers.VerifyOTP)
			otp.POST("/resend", controllers.ResendOTP)
		}
	}
}
