package models

type EmailReset struct {
	Email       string `json:"email"`
	OTPCode     string `json:"otp_code"`
}
