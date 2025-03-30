package dto

type CreateUser struct {
	Name     string `binding:"required"       json:"name"`
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required,min=8" json:"password"`
}

type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
