package models

type PasswordReset struct {
	Email       string `json:"email"`
	Password     string `json:"password"`
}