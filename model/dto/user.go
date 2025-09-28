package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=8"`
	Role        string  `json:"role" binding:"required,oneof=student business admin"`
	CompanyName *string `json:"company_name" binding:"required_if=Role business"`
	University  *string `json:"university" binding:"required_if=Role student"`
}
