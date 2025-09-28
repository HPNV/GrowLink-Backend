package db

type Project struct {
	UUID         string `db:"uuid"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	Status       string `db:"status"`
	Duration     int    `db:"duration"`
	Timeline     string `db:"timeline"`
	Deliverables string `db:"deliverables"`
	CreatedBy    string `db:"created_by"`
	CreatedAt    string `db:"created_at"`
}
