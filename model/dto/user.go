package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Name        string  `json:"name" binding:"required"`
	Password    string  `json:"password" binding:"required,min=8"`
	Role        string  `json:"role" binding:"required,oneof=student business admin"`
	CompanyName *string `json:"company_name" binding:"required_if=Role business"`
	University  *string `json:"university" binding:"required_if=Role student"`
}

type UserResponse struct {
	UUID      string `json:"uuid"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserDetailResponse struct {
	UUID        string   `json:"uuid"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	CompanyName *string  `json:"company_name,omitempty"`
	University  *string  `json:"university,omitempty"`
	CreatedAt   string   `json:"created_at"`
	Skills      []string `json:"skills,omitempty"`
}

type StudentListRequest struct {
	Name       *string `json:"name"`
	University *string `json:"university"`
	Skill      *string `json:"skill"`
	Page       int     `json:"page" binding:"required"`
	Limit      int     `json:"limit" binding:"required"`
}

type StudentListResponse struct {
	Students   []*StudentDetailResponse `json:"students"`
	TotalCount int                      `json:"total_count"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

type StudentDetailResponse struct {
	UUID       string   `json:"uuid"`
	UserUUID   string   `json:"user_uuid"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	University string   `json:"university"`
	Skills     []string `json:"skills"`
	CreatedAt  string   `json:"created_at"`
}
