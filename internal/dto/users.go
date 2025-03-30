package dto

type CreateUser struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
