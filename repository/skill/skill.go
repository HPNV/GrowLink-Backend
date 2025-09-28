package skill

import (
	"context"

	"github.com/HPNV/growlink-backend/model/db"
	"github.com/jmoiron/sqlx"
)

type ISkill interface {
	CreateSkill(ctx context.Context, tx *sqlx.Tx, name string) (string, error)
	Create(tx *sqlx.Tx, skill *db.Skill) error
	GetByUUID(uuid string) (*db.Skill, error)
	Update(tx *sqlx.Tx, skill *db.Skill) error
	Delete(tx *sqlx.Tx, uuid string) error
	GetAll() ([]*db.Skill, error)
	GetByProjectUUID(projectUUID string) ([]*db.Skill, error)
}

type Skill struct {
	db *sqlx.DB
}

func NewSkill(db *sqlx.DB) ISkill {
	return &Skill{
		db: db,
	}
}

func (s *Skill) CreateSkill(ctx context.Context, tx *sqlx.Tx, name string) (string, error) {
	var uuid string
	err := tx.QueryRowContext(ctx, "INSERT INTO skills (name) VALUES ($1) RETURNING uuid", name).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (s *Skill) Create(tx *sqlx.Tx, skill *db.Skill) error {
	query := `
		INSERT INTO skills (name, description)
		VALUES ($1, $2)
		RETURNING uuid, created_at
	`
	return tx.QueryRow(query, skill.Name, skill.Description).Scan(&skill.UUID, &skill.CreatedAt)
}

func (s *Skill) GetByUUID(uuid string) (*db.Skill, error) {
	skill := &db.Skill{}
	query := `SELECT uuid, name, description, created_at FROM skills WHERE uuid = $1`
	err := s.db.Get(skill, query, uuid)
	return skill, err
}

func (s *Skill) Update(tx *sqlx.Tx, skill *db.Skill) error {
	query := `
		UPDATE skills 
		SET name = $1, description = $2
		WHERE uuid = $3
	`
	_, err := tx.Exec(query, skill.Name, skill.Description, skill.UUID)
	return err
}

func (s *Skill) Delete(tx *sqlx.Tx, uuid string) error {
	query := `DELETE FROM skills WHERE uuid = $1`
	_, err := tx.Exec(query, uuid)
	return err
}

func (s *Skill) GetAll() ([]*db.Skill, error) {
	var skills []*db.Skill
	query := `SELECT uuid, name, description, created_at FROM skills ORDER BY name`
	err := s.db.Select(&skills, query)
	return skills, err
}

func (s *Skill) GetByProjectUUID(projectUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	query := `
		SELECT s.uuid, s.name, s.description, s.created_at
		FROM skills s
		JOIN project_skills ps ON s.uuid = ps.skill_uuid
		WHERE ps.project_uuid = $1
	`
	err := s.db.Select(&skills, query, projectUUID)
	return skills, err
}
