package dto

type StudentRequest struct {
	University string `json:"university" binding:"required"`
}

type StudentResponse struct {
	UUID       string `json:"uuid"`
	UserUUID   string `json:"user_uuid"`
	University string `json:"university"`
}
