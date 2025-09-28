package project

import (
	"github.com/HPNV/growlink-backend/model/db"
	"github.com/jmoiron/sqlx"
)

type IProject interface {
	Create(tx *sqlx.Tx, project *db.Project) error
	GetByUUID(uuid string) (*db.Project, error)
	Update(tx *sqlx.Tx, project *db.Project) error
	Delete(tx *sqlx.Tx, uuid string) error
	GetAll() ([]*db.Project, error)
	GetByBusinessUUID(businessUUID string) ([]*db.Project, error)
	AddSkill(tx *sqlx.Tx, projectUUID, skillUUID string) error
	RemoveSkill(tx *sqlx.Tx, projectUUID, skillUUID string) error
	GetSkills(projectUUID string) ([]*db.Skill, error)
	AddStudent(tx *sqlx.Tx, projectUUID, studentUUID string) error
	RemoveStudent(tx *sqlx.Tx, projectUUID, studentUUID string) error
	GetStudents(projectUUID string) ([]*db.Student, error)
}

type Project struct {
	db *sqlx.DB
}

func NewProject(db *sqlx.DB) IProject {
	return &Project{
		db: db,
	}
}

func (p *Project) Create(tx *sqlx.Tx, project *db.Project) error {
	query := `
		INSERT INTO projects (name, description, duration, timeline, deliverables, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid, created_at
	`

	err := tx.QueryRow(query, project.Name, project.Description, project.Duration, project.Timeline, project.Deliverables, project.CreatedBy).Scan(&project.UUID, &project.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) GetByUUID(uuid string) (*db.Project, error) {
	project := &db.Project{}
	query := `SELECT uuid, name, description, status, created_by, created_at FROM projects WHERE uuid = $1`
	err := p.db.Get(project, query, uuid)
	return project, err
}

func (p *Project) Update(tx *sqlx.Tx, project *db.Project) error {
	query := `
		UPDATE projects 
		SET name = $1, description = $2, status = $3
		WHERE uuid = $4
	`
	_, err := tx.Exec(query, project.Name, project.Description, project.Status, project.UUID)
	return err
}

func (p *Project) Delete(tx *sqlx.Tx, uuid string) error {
	query := `DELETE FROM projects WHERE uuid = $1`
	_, err := tx.Exec(query, uuid)
	return err
}

func (p *Project) GetAll() ([]*db.Project, error) {
	var projects []*db.Project
	query := `SELECT uuid, name, description, status, created_by, created_at FROM projects ORDER BY created_at DESC`
	err := p.db.Select(&projects, query)
	return projects, err
}

func (p *Project) GetByBusinessUUID(businessUUID string) ([]*db.Project, error) {
	var projects []*db.Project
	query := `SELECT uuid, name, description, status, created_by, created_at FROM projects WHERE created_by = $1 ORDER BY created_at DESC`
	err := p.db.Select(&projects, query, businessUUID)
	return projects, err
}

func (p *Project) AddSkill(tx *sqlx.Tx, projectUUID, skillUUID string) error {
	query := `
		INSERT INTO project_skills (project_uuid, skill_uuid)
		VALUES ($1, $2)
		ON CONFLICT (project_uuid, skill_uuid) DO NOTHING
	`
	_, err := tx.Exec(query, projectUUID, skillUUID)
	return err
}

func (p *Project) RemoveSkill(tx *sqlx.Tx, projectUUID, skillUUID string) error {
	query := `DELETE FROM project_skills WHERE project_uuid = $1 AND skill_uuid = $2`
	_, err := tx.Exec(query, projectUUID, skillUUID)
	return err
}

func (p *Project) GetSkills(projectUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	query := `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		INNER JOIN project_skills ps ON s.uuid = ps.skill_uuid
		WHERE ps.project_uuid = $1
		ORDER BY s.name
	`
	err := p.db.Select(&skills, query, projectUUID)
	return skills, err
}

func (p *Project) AddStudent(tx *sqlx.Tx, projectUUID, studentUUID string) error {
	query := `
		INSERT INTO student_projects (student_uuid, project_uuid)
		VALUES ($1, $2)
		ON CONFLICT (student_uuid, project_uuid) DO NOTHING
	`
	_, err := tx.Exec(query, studentUUID, projectUUID)
	return err
}

func (p *Project) RemoveStudent(tx *sqlx.Tx, projectUUID, studentUUID string) error {
	query := `DELETE FROM student_projects WHERE project_uuid = $1 AND student_uuid = $2`
	_, err := tx.Exec(query, projectUUID, studentUUID)
	return err
}

func (p *Project) GetStudents(projectUUID string) ([]*db.Student, error) {
	var students []*db.Student
	query := `
		SELECT s.uuid, s.user_uuid, s.university
		FROM students s
		INNER JOIN student_projects sp ON s.uuid = sp.student_uuid
		WHERE sp.project_uuid = $1
		ORDER BY s.university
	`
	err := p.db.Select(&students, query, projectUUID)
	return students, err
}
