package dto

type ProjectRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	Duration     int      `json:"duration" binding:"required"`
	Timeline     string   `json:"timeline" binding:"required,oneof=day week month year"`
	Deliverables string   `json:"deliverables" binding:"required"`
	Skills       []string `json:"skills"`
}

type ProjectResponse struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Status       string   `json:"status"`
	Duration     int      `json:"duration"`
	Timeline     string   `json:"timeline"`
	Deliverables string   `json:"deliverables"`
	Skills       []string `json:"skills"`
	CreatedBy    string   `json:"created_by"`
	CreatedAt    string   `json:"created_at"`
}

type ProjectUpdateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Status       string `json:"status" binding:"omitempty,oneof=open in_progress completed"`
	Duration     *int   `json:"duration"`
	Timeline     string `json:"timeline" binding:"omitempty,oneof=day week month year"`
	Deliverables string `json:"deliverables"`
}

type ProjectListRequest struct {
	Skill  *string `json:"skill"`
	Budget *int    `json:"budget"`
	Search *string `json:"search"`
	Page   int     `json:"page" binding:"required"`
	Limit  int     `json:"limit" binding:"required"`
}

type ProjectListResponse struct {
	Projects   []*ProjectResponse `json:"projects"`
	TotalCount int                `json:"total_count"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}
