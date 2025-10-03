package project

import (
	"github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/jmoiron/sqlx"
)

type IProject interface {
	Create(businessUUID string, req *dto.ProjectRequest) (*dto.ProjectResponse, error)
	GetByUUID(uuid string) (*dto.ProjectResponse, error)
	Update(uuid string, req *dto.ProjectUpdateRequest) (*dto.ProjectResponse, error)
	Delete(uuid string) error
	GetAll() ([]*dto.ProjectResponse, error)
	GetAllList(req *dto.ProjectListRequest) (*dto.ProjectListResponse, error)
	GetByBusinessUUID(businessUUID string) ([]*dto.ProjectResponse, error)
	AddSkill(projectUUID, skillName string) error
	RemoveSkill(projectUUID, skillName string) error
	GetSkills(projectUUID string) ([]*dto.SkillResponse, error)
	AddStudent(projectUUID, studentUUID string) error
	RemoveStudent(projectUUID, studentUUID string) error
	GetStudents(projectUUID string) ([]*dto.StudentResponse, error)
}

type Project struct {
	repo repository.IRegistry
}

func NewProject(repo repository.IRegistry) IProject {
	return &Project{
		repo: repo,
	}
}

func (p *Project) Create(businessUUID string, req *dto.ProjectRequest) (*dto.ProjectResponse, error) {
	project := &db.Project{
		Name:         req.Name,
		Description:  req.Description,
		Duration:     req.Duration,
		Timeline:     req.Timeline,
		Deliverables: req.Deliverables,
		Status:       "open",
		CreatedBy:    businessUUID,
	}

	err := p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		if err := p.repo.GetProject().Create(tx, project); err != nil {
			return err
		}

		for _, skillName := range req.Skills {
			// Get skill by name to get its UUID
			skill, err := p.repo.GetSkill().GetByName(skillName)
			if err != nil {
				return err
			}
			if err := p.repo.GetProject().AddSkill(tx, project.UUID, skill.UUID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.ProjectResponse{
		UUID:         project.UUID,
		Name:         project.Name,
		Description:  project.Description,
		Status:       project.Status,
		Duration:     project.Duration,
		Timeline:     project.Timeline,
		Deliverables: project.Deliverables,
		Skills:       req.Skills,
		CreatedBy:    project.CreatedBy,
		CreatedAt:    project.CreatedAt,
	}, nil
}

func (p *Project) GetByUUID(uuid string) (*dto.ProjectResponse, error) {
	project, err := p.repo.GetProject().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	skills, err := p.repo.GetSkill().GetByProjectUUID(project.UUID)
	if err != nil {
		return nil, err
	}

	var skillNames []string
	for _, skill := range skills {
		skillNames = append(skillNames, skill.Name)
	}

	return &dto.ProjectResponse{
		UUID:         project.UUID,
		Name:         project.Name,
		Description:  project.Description,
		Status:       project.Status,
		Deliverables: project.Deliverables,
		Duration:     project.Duration,
		Timeline:     project.Timeline,
		Skills:       skillNames,
		CreatedBy:    project.CreatedBy,
		CreatedAt:    project.CreatedAt,
	}, nil
}

func (p *Project) Update(uuid string, req *dto.ProjectUpdateRequest) (*dto.ProjectResponse, error) {
	// First get existing project
	existing, err := p.repo.GetProject().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	if req.Duration != nil {
		existing.Duration = *req.Duration
	}
	if req.Timeline != "" {
		existing.Timeline = req.Timeline
	}
	if req.Deliverables != "" {
		existing.Deliverables = req.Deliverables
	}

	err = p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return p.repo.GetProject().Update(tx, existing)
	})

	if err != nil {
		return nil, err
	}

	// Get skills for the response
	skills, err := p.repo.GetSkill().GetByProjectUUID(existing.UUID)
	if err != nil {
		return nil, err
	}

	var skillNames []string
	for _, skill := range skills {
		skillNames = append(skillNames, skill.Name)
	}

	return &dto.ProjectResponse{
		UUID:         existing.UUID,
		Name:         existing.Name,
		Description:  existing.Description,
		Status:       existing.Status,
		Duration:     existing.Duration,
		Timeline:     existing.Timeline,
		Deliverables: existing.Deliverables,
		Skills:       skillNames,
		CreatedBy:    existing.CreatedBy,
		CreatedAt:    existing.CreatedAt,
	}, nil
}

func (p *Project) Delete(uuid string) error {
	return p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return p.repo.GetProject().Delete(tx, uuid)
	})
}

func (p *Project) GetAll() ([]*dto.ProjectResponse, error) {
	projects, err := p.repo.GetProject().GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*dto.ProjectResponse
	for _, project := range projects {
		responses = append(responses, &dto.ProjectResponse{
			UUID:         project.UUID,
			Name:         project.Name,
			Description:  project.Description,
			Status:       project.Status,
			Duration:     project.Duration,
			Timeline:     project.Timeline,
			Deliverables: project.Deliverables,
			CreatedBy:    project.CreatedBy,
			CreatedAt:    project.CreatedAt,
		})
	}

	for _, project := range responses {
		skills, err := p.repo.GetSkill().GetByProjectUUID(project.UUID)
		if err != nil {
			return nil, err
		}
		for _, skill := range skills {
			project.Skills = append(project.Skills, skill.Name)
		}
	}

	return responses, nil
}

func (p *Project) GetAllList(req *dto.ProjectListRequest) (*dto.ProjectListResponse, error) {
	projects, totalCount, err := p.repo.GetProject().GetAllList(req)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ProjectResponse
	for _, project := range projects {
		responses = append(responses, &dto.ProjectResponse{
			UUID:         project.UUID,
			Name:         project.Name,
			Description:  project.Description,
			Status:       project.Status,
			Duration:     project.Duration,
			Timeline:     project.Timeline,
			Deliverables: project.Deliverables,
			CreatedBy:    project.CreatedBy,
			CreatedAt:    project.CreatedAt,
		})
	}

	for _, project := range responses {
		skills, err := p.repo.GetSkill().GetByProjectUUID(project.UUID)
		if err != nil {
			return nil, err
		}
		for _, skill := range skills {
			project.Skills = append(project.Skills, skill.Name)
		}
	}

	totalPages := totalCount / req.Limit
	if totalCount%req.Limit != 0 {
		totalPages++
	}

	return &dto.ProjectListResponse{
		Projects:   responses,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (p *Project) GetByBusinessUUID(businessUUID string) ([]*dto.ProjectResponse, error) {
	projects, err := p.repo.GetProject().GetByBusinessUUID(businessUUID)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ProjectResponse
	for _, project := range projects {
		responses = append(responses, &dto.ProjectResponse{
			UUID:         project.UUID,
			Name:         project.Name,
			Description:  project.Description,
			Status:       project.Status,
			Duration:     project.Duration,
			Timeline:     project.Timeline,
			Deliverables: project.Deliverables,
			CreatedBy:    project.CreatedBy,
			CreatedAt:    project.CreatedAt,
		})
	}

	// Load skills for each project
	for _, project := range responses {
		skills, err := p.repo.GetSkill().GetByProjectUUID(project.UUID)
		if err != nil {
			return nil, err
		}
		for _, skill := range skills {
			project.Skills = append(project.Skills, skill.Name)
		}
	}

	return responses, nil
}

func (p *Project) AddSkill(projectUUID, skillName string) error {
	// First, get the skill by name to get its UUID
	skill, err := p.repo.GetSkill().GetByName(skillName)
	if err != nil {
		return err
	}

	return p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return p.repo.GetProject().AddSkill(tx, projectUUID, skill.UUID)
	})
}

func (p *Project) RemoveSkill(projectUUID, skillName string) error {
	// First, get the skill by name to get its UUID
	skill, err := p.repo.GetSkill().GetByName(skillName)
	if err != nil {
		return err
	}

	return p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return p.repo.GetProject().RemoveSkill(tx, projectUUID, skill.UUID)
	})
}

func (p *Project) GetSkills(projectUUID string) ([]*dto.SkillResponse, error) {
	skills, err := p.repo.GetProject().GetSkills(projectUUID)
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

func (p *Project) AddStudent(projectUUID, studentUUID string) error {
	return p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return p.repo.GetProject().AddStudent(tx, projectUUID, studentUUID)
	})
}

func (p *Project) RemoveStudent(projectUUID, studentUUID string) error {
	return p.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return p.repo.GetProject().RemoveStudent(tx, projectUUID, studentUUID)
	})
}

func (p *Project) GetStudents(projectUUID string) ([]*dto.StudentResponse, error) {
	students, err := p.repo.GetProject().GetStudents(projectUUID)
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
