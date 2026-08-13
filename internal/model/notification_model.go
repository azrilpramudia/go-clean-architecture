package model

type SendVerificationEmailRequest struct {
	Email string `json:"email"`
	Code string `json:"code"`
}