package dto

type BusinessRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
}

type BusinessResponse struct {
	UUID        string `json:"uuid"`
	UserUUID    string `json:"user_uuid"`
	CompanyName string `json:"company_name"`
}
