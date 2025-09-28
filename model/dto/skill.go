package dto

type SkillRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type SkillResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
