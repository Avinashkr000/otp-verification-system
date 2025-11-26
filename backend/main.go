package main

import (
	"fmt"
	"log"
	"otp-backend/config"
	"otp-backend/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	fmt.Println("\n🔧 Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  Warning: .env file not found, using system environment variables")
	} else {
		fmt.Println("✅ .env file loaded successfully")
		
		// Show Twilio configuration status
		if os.Getenv("TWILIO_ACCOUNT_SID") != "" {
			fmt.Println("✅ Twilio SMS configured")
			fmt.Printf("   - Account SID: %s...\n", os.Getenv("TWILIO_ACCOUNT_SID")[:10])
			fmt.Printf("   - Phone Number: %s\n", os.Getenv("TWILIO_PHONE_NUMBER"))
		} else {
			fmt.Println("⚠️  Twilio not configured - SMS sending disabled")
		}
	}
	fmt.Println()

	// Initialize database connection
	config.ConnectDatabase()

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "OTP Verification API is running",
		})
	})

	// Register routes
	routes.RegisterOTPRoutes(router)

	// Start server
	log.Println("\n🚀 Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
