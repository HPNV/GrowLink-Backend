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
	GetByName(name string) (*db.Skill, error)
	Update(tx *sqlx.Tx, skill *db.Skill) error
	Delete(tx *sqlx.Tx, uuid string) error
	GetAll() ([]*db.Skill, error)
	GetByProjectUUID(projectUUID string) ([]*db.Skill, error)
	GetByStudentUUID(studentUUID string) ([]*db.Skill, error)
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
	err := tx.QueryRowContext(ctx, CreateSkillQuery, name).Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (s *Skill) Create(tx *sqlx.Tx, skill *db.Skill) error {
	return tx.QueryRow(CreateQuery, skill.Name, skill.Description).Scan(&skill.UUID, &skill.CreatedAt)
}

func (s *Skill) GetByUUID(uuid string) (*db.Skill, error) {
	skill := &db.Skill{}
	err := s.db.Get(skill, GetByUUIDQuery, uuid)
	return skill, err
}

func (s *Skill) GetByName(name string) (*db.Skill, error) {
	skill := &db.Skill{}
	err := s.db.Get(skill, GetByNameQuery, name)
	return skill, err
}

func (s *Skill) Update(tx *sqlx.Tx, skill *db.Skill) error {
	_, err := tx.Exec(UpdateQuery, skill.Name, skill.Description, skill.UUID)
	return err
}

func (s *Skill) Delete(tx *sqlx.Tx, uuid string) error {
	_, err := tx.Exec(DeleteQuery, uuid)
	return err
}

func (s *Skill) GetAll() ([]*db.Skill, error) {
	var skills []*db.Skill
	err := s.db.Select(&skills, GetAllQuery)
	return skills, err
}

func (s *Skill) GetByProjectUUID(projectUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	err := s.db.Select(&skills, GetByProjectUUIDQuery, projectUUID)
	return skills, err
}

func (s *Skill) GetByStudentUUID(studentUUID string) ([]*db.Skill, error) {
	var skills []*db.Skill
	err := s.db.Select(&skills, GetByStudentUUIDQuery, studentUUID)
	return skills, err
}
