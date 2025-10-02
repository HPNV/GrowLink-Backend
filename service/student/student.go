package student

import (
	"github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/jmoiron/sqlx"
)

type IStudent interface {
	Create(userUUID string, req *dto.StudentRequest) (*dto.StudentResponse, error)
	GetByUUID(uuid string) (*dto.StudentResponse, error)
	GetByUserUUID(userUUID string) (*dto.StudentResponse, error)
	Update(uuid string, req *dto.StudentRequest) (*dto.StudentResponse, error)
	Delete(uuid string) error
	GetAll() ([]*dto.StudentResponse, error)
	AddSkill(studentUUID, skillName string) error
	RemoveSkill(studentUUID, skillName string) error
	GetSkills(studentUUID string) ([]*dto.SkillResponse, error)
}

type Student struct {
	repo repository.IRegistry
}

func NewStudent(repo repository.IRegistry) IStudent {
	return &Student{
		repo: repo,
	}
}

func (s *Student) Create(userUUID string, req *dto.StudentRequest) (*dto.StudentResponse, error) {
	student := &db.Student{
		UserUUID:   userUUID,
		University: req.University,
	}

	err := s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetStudent().Create(tx, student)
	})

	if err != nil {
		return nil, err
	}

	return &dto.StudentResponse{
		UUID:       student.UUID,
		UserUUID:   student.UserUUID,
		University: student.University,
	}, nil
}

func (s *Student) GetByUUID(uuid string) (*dto.StudentResponse, error) {
	student, err := s.repo.GetStudent().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	return &dto.StudentResponse{
		UUID:       student.UUID,
		UserUUID:   student.UserUUID,
		University: student.University,
	}, nil
}

func (s *Student) GetByUserUUID(userUUID string) (*dto.StudentResponse, error) {
	student, err := s.repo.GetStudent().GetByUserUUID(userUUID)
	if err != nil {
		return nil, err
	}

	return &dto.StudentResponse{
		UUID:       student.UUID,
		UserUUID:   student.UserUUID,
		University: student.University,
	}, nil
}

func (s *Student) Update(uuid string, req *dto.StudentRequest) (*dto.StudentResponse, error) {
	// First get existing student
	existing, err := s.repo.GetStudent().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.University = req.University

	err = s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetStudent().Update(tx, existing)
	})

	if err != nil {
		return nil, err
	}

	return &dto.StudentResponse{
		UUID:       existing.UUID,
		UserUUID:   existing.UserUUID,
		University: existing.University,
	}, nil
}

func (s *Student) Delete(uuid string) error {
	return s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetStudent().Delete(tx, uuid)
	})
}

func (s *Student) GetAll() ([]*dto.StudentResponse, error) {
	students, err := s.repo.GetStudent().GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*dto.StudentResponse
	for _, student := range students {
		responses = append(responses, &dto.StudentResponse{
			UUID:       student.UUID,
			UserUUID:   student.UserUUID,
			University: student.University,
		})
	}

	return responses, nil
}

func (s *Student) AddSkill(studentUUID, skillName string) error {
	// First, get the skill by name to get its UUID
	skill, err := s.repo.GetSkill().GetByName(skillName)
	if err != nil {
		return err
	}

	return s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetStudent().AddSkill(tx, studentUUID, skill.UUID)
	})
}

func (s *Student) RemoveSkill(studentUUID, skillName string) error {
	// First, get the skill by name to get its UUID
	skill, err := s.repo.GetSkill().GetByName(skillName)
	if err != nil {
		return err
	}

	return s.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return s.repo.GetStudent().RemoveSkill(tx, studentUUID, skill.UUID)
	})
}

func (s *Student) GetSkills(studentUUID string) ([]*dto.SkillResponse, error) {
	skills, err := s.repo.GetStudent().GetSkills(studentUUID)
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
