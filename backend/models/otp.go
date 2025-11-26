package models

import (
	"time"
)

type OTP struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email        string     `gorm:"type:varchar(255)" json:"email"`
	Phone        string     `gorm:"type:varchar(20)" json:"phone"`
	OTPCode      string     `gorm:"type:varchar(6);not null" json:"otp_code"`
	IsVerified   bool       `gorm:"default:false" json:"is_verified"`
	AttemptCount int        `gorm:"default:0" json:"attempt_count"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	VerifiedAt   *time.Time `json:"verified_at"`
}

type User struct {
	ID               string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email            string    `gorm:"type:varchar(255);unique" json:"email"`
	Phone            string    `gorm:"type:varchar(20);unique" json:"phone"`
	IsEmailVerified  bool      `gorm:"default:false" json:"is_email_verified"`
	IsPhoneVerified  bool      `gorm:"default:false" json:"is_phone_verified"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
