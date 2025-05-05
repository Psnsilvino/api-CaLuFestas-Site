package models

type EmailReset struct {
	Email string `json:"email"`
	NewPassword string `json:"newPassword"`
	OTPCode    string             `json:"otp_code"`
}
