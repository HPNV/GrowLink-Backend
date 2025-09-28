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
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=open in_progress completed"`
}
