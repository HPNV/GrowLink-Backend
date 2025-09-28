package skill

import (
	"context"

	"github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/jmoiron/sqlx"
)

type ISkill interface {
	CreateSkill(name string) (string, error)
	Create(req *dto.SkillRequest) (*dto.SkillResponse, error)
	GetByUUID(uuid string) (*dto.SkillResponse, error)
	Update(uuid string, req *dto.SkillRequest) (*dto.SkillResponse, error)
	Delete(uuid string) error
	GetAll() ([]*dto.SkillResponse, error)
}

type Skill struct {
	repo repository.IRegistry
}

func NewSkill(repo repository.IRegistry) ISkill {
	return &Skill{
		repo: repo,
	}
}

func (s *Skill) CreateSkill(name string) (string, error) {
	var skillID string
	err := s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		var err error
		skillID, err = s.repo.GetSkill().CreateSkill(context.Background(), tx, name)
		return err
	})
	if err != nil {
		return "", err
	}
	return skillID, nil
}

func (s *Skill) Create(req *dto.SkillRequest) (*dto.SkillResponse, error) {
	skill := &db.Skill{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetSkill().Create(tx, skill)
	})

	if err != nil {
		return nil, err
	}

	return &dto.SkillResponse{
		UUID:        skill.UUID,
		Name:        skill.Name,
		Description: skill.Description,
		CreatedAt:   skill.CreatedAt,
	}, nil
}

func (s *Skill) GetByUUID(uuid string) (*dto.SkillResponse, error) {
	skill, err := s.repo.GetSkill().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	return &dto.SkillResponse{
		UUID:        skill.UUID,
		Name:        skill.Name,
		Description: skill.Description,
		CreatedAt:   skill.CreatedAt,
	}, nil
}

func (s *Skill) Update(uuid string, req *dto.SkillRequest) (*dto.SkillResponse, error) {
	// First get existing skill
	existing, err := s.repo.GetSkill().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.Name = req.Name
	existing.Description = req.Description

	err = s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetSkill().Update(tx, existing)
	})

	if err != nil {
		return nil, err
	}

	return &dto.SkillResponse{
		UUID:        existing.UUID,
		Name:        existing.Name,
		Description: existing.Description,
		CreatedAt:   existing.CreatedAt,
	}, nil
}

func (s *Skill) Delete(uuid string) error {
	return s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetSkill().Delete(tx, uuid)
	})
}

func (s *Skill) GetAll() ([]*dto.SkillResponse, error) {
	skills, err := s.repo.GetSkill().GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*dto.SkillResponse
	for _, skill := range skills {
		responses = append(responses, &dto.SkillResponse{
			UUID:        skill.UUID,
			Name:        skill.Name,
			Description: skill.Description,
			CreatedAt:   skill.CreatedAt,
		})
	}

	return responses, nil
}
