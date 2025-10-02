package project

import (
	"fmt"
	"strings"

	"github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/jmoiron/sqlx"
)

type IProject interface {
	Create(tx *sqlx.Tx, project *db.Project) error
	GetByUUID(uuid string) (*db.Project, error)
	Update(tx *sqlx.Tx, project *db.Project) error
	Delete(tx *sqlx.Tx, uuid string) error
	GetAll() ([]*db.Project, error)
	GetAllList(queryParam *dto.ProjectListRequest) ([]*db.Project, int, error)
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
	err := tx.QueryRow(CreateQuery, project.Name, project.Description, project.Duration, project.Timeline, project.Deliverables, project.CreatedBy).Scan(&project.UUID, &project.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) GetByUUID(uuid string) (*db.Project, error) {
	project := &db.Project{}
	err := p.db.Get(project, GetByUUIDQuery, uuid)
	return project, err
}

func (p *Project) Update(tx *sqlx.Tx, project *db.Project) error {
	_, err := tx.Exec(UpdateQuery, project.Name, project.Description, project.Status, project.UUID)
	return err
}

func (p *Project) Delete(tx *sqlx.Tx, uuid string) error {
	_, err := tx.Exec(DeleteQuery, uuid)
	return err
}

func (p *Project) GetAll() ([]*db.Project, error) {
	var projects []*db.Project
	err := p.db.Select(&projects, GetAllQuery)
	return projects, err
}

func (p *Project) GetByBusinessUUID(businessUUID string) ([]*db.Project, error) {
	var projects []*db.Project
	err := p.db.Select(&projects, GetByBusinessUUIDQuery, businessUUID)
	return projects, err
}

func (p *Project) AddSkill(tx *sqlx.Tx, projectUUID, skillUUID string) error {
	_, err := tx.Exec(AddSkillQuery, projectUUID, skillUUID)
	return err
}

func (p *Project) RemoveSkill(tx *sqlx.Tx, projectUUID, skillUUID string) error {
	_, err := tx.Exec(RemoveSkillQuery, projectUUID, skillUUID)
	return err
}

func (p *Project) GetSkills(projectUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	err := p.db.Select(&skills, GetSkillsQuery, projectUUID)
	return skills, err
}

func (p *Project) AddStudent(tx *sqlx.Tx, projectUUID, studentUUID string) error {
	_, err := tx.Exec(AddStudentQuery, studentUUID, projectUUID)
	return err
}

func (p *Project) RemoveStudent(tx *sqlx.Tx, projectUUID, studentUUID string) error {
	_, err := tx.Exec(RemoveStudentQuery, projectUUID, studentUUID)
	return err
}

func (p *Project) GetStudents(projectUUID string) ([]*db.Student, error) {
	var students []*db.Student
	err := p.db.Select(&students, GetStudentsQuery, projectUUID)
	return students, err
}

func (p *Project) GetAllList(queryParam *dto.ProjectListRequest) ([]*db.Project, int, error) {
	var projects []*db.Project
	var args []interface{}
	var whereConditions []string
	argIndex := 1

	// Build WHERE conditions
	if queryParam.Skill != nil && *queryParam.Skill != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(`EXISTS (
			SELECT 1 FROM project_skills ps 
			JOIN skills s ON ps.skill_uuid = s.uuid 
			WHERE ps.project_uuid = p.uuid AND s.name ILIKE $%d
		)`, argIndex))
		args = append(args, "%"+*queryParam.Skill+"%")
		argIndex++
	}

	if queryParam.Budget != nil && *queryParam.Budget > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("p.duration <= $%d", argIndex))
		args = append(args, *queryParam.Budget)
		argIndex++
	}

	if queryParam.Search != nil && *queryParam.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(`(
			p.name ILIKE $%d OR 
			p.description ILIKE $%d OR 
			p.deliverables ILIKE $%d
		)`, argIndex, argIndex, argIndex))
		args = append(args, "%"+*queryParam.Search+"%")
		argIndex++
	}

	// Build the complete query
	query := GetAllListQuery
	countQuery := GetAllListCountQuery

	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// Get total count
	var totalCount int
	err := p.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add ORDER BY and pagination
	query += " ORDER BY p.created_at DESC"

	if queryParam.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, queryParam.Limit)
		argIndex++
	}

	if queryParam.Page > 0 {
		offset := (queryParam.Page - 1) * queryParam.Limit
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	// Execute the query
	err = p.db.Select(&projects, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return projects, totalCount, nil
}
