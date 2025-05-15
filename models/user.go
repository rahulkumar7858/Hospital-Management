package models

type User struct {
	ID                   int    `json:"id"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	Role                 string `json:"role"`
	ForgotPasswordToken  string `json:"forgot_password_token"`
	ForgotPasswordExpiry string `json:"forgot_password_expiry"`
}
