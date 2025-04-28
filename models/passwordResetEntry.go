package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordResetEntry struct {
	ID         primitive.ObjectID `json:"id"`
	Email      string             `json:"email"`
	OTPCode    string             `json:"otp_code"`
	ExpiresAt  time.Time          `json:"expires_at"`
	IsVerified bool               `json:"is_verified"`
	CreatedAt  time.Time          `json:"created_at"`
}
