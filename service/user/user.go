package user

import (
	"context"
	"errors"

	modelDB "github.com/HPNV/growlink-backend/model/db"
	modelDTO "github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/jmoiron/sqlx"
)

type IUser interface {
	Login(ctx context.Context, email, password string) (*modelDB.User, error)
	Register(ctx context.Context, request modelDTO.RegisterRequest) (*modelDB.User, error)
	GetAll() ([]*modelDTO.UserResponse, error)
	GetDetail(ctx context.Context, uuid string) (*modelDTO.UserDetailResponse, error)
	GetStudentList(req *modelDTO.StudentListRequest) (*modelDTO.StudentListResponse, error)
}

type User struct {
	repo repository.IRegistry
}

func NewUser(repo repository.IRegistry) IUser {
	return &User{
		repo: repo,
	}
}

func (u *User) Login(ctx context.Context, email, password string) (*modelDB.User, error) {
	user, err := u.repo.GetUser().Login(
		ctx,
		email,
		password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Register(ctx context.Context, request modelDTO.RegisterRequest) (*modelDB.User, error) {
	user := &modelDB.User{
		Email: request.Email,
		Name:  request.Name,
		Role:  request.Role,
	}
	txFunc := func(tx *sqlx.Tx) error {
		userResult, err := u.repo.GetUser().Register(ctx, tx, user, request.Password)
		if err != nil {
			return err
		}

		switch user.Role {
		case "student":
			if request.University == nil {
				return errors.New("university is required for student role")
			}
			student := &modelDB.Student{
				UserUUID:   userResult.UUID,
				University: *request.University,
			}
			err = u.repo.GetStudent().Create(tx, student)
			if err != nil {
				return err
			}
		case "business":
			if request.CompanyName == nil {
				return errors.New("company_name is required for business role")
			}
			business := &modelDB.Business{
				UserUUID:    userResult.UUID,
				CompanyName: *request.CompanyName,
			}
			err = u.repo.GetBusiness().Create(tx, business)
			if err != nil {
				return err
			}
		}

		user = userResult

		return nil
	}

	err := u.repo.WithTransaction(txFunc)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetAll() ([]*modelDTO.UserResponse, error) {
	users, err := u.repo.GetUser().GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*modelDTO.UserResponse
	for _, user := range users {
		responses = append(responses, &modelDTO.UserResponse{
			UUID:      user.UUID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		})
	}

	return responses, nil
}

func (u *User) GetDetail(ctx context.Context, uuid string) (*modelDTO.UserDetailResponse, error) {
	user, err := u.repo.GetUser().GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	userDTO := modelDTO.UserDetailResponse{
		UUID:      user.UUID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	switch user.Role {
	case "student":
		student, err := u.repo.GetStudent().GetByUserUUID(user.UUID)
		if err != nil {
			return nil, err
		}
		userDTO.University = &student.University

		skills, err := u.repo.GetSkill().GetByStudentUUID(user.UUID)
		if err != nil {
			return nil, err
		}

		for _, skill := range skills {
			userDTO.Skills = append(userDTO.Skills, skill.Name)
		}
	case "business":
		business, err := u.repo.GetBusiness().GetByUserUUID(user.UUID)
		if err != nil {
			return nil, err
		}
		userDTO.CompanyName = &business.CompanyName
	}

	return &userDTO, nil
}

func (u *User) GetStudentList(req *modelDTO.StudentListRequest) (*modelDTO.StudentListResponse, error) {
	users, students, totalCount, err := u.repo.GetUser().GetStudentList(req)
	if err != nil {
		return nil, err
	}

	var responses []*modelDTO.StudentDetailResponse

	// Create a map for quick lookup of students by user_uuid
	studentMap := make(map[string]*modelDB.Student)
	for _, student := range students {
		studentMap[student.UserUUID] = student
	}

	for _, user := range users {
		student := studentMap[user.UUID]
		if student == nil {
			continue // Skip if no student record found
		}

		studentResponse := &modelDTO.StudentDetailResponse{
			UUID:       student.UUID,
			UserUUID:   user.UUID,
			Email:      user.Email,
			Name:       user.Name,
			University: student.University,
			CreatedAt:  user.CreatedAt,
		}

		// Get skills for this student
		skills, err := u.repo.GetSkill().GetByStudentUUID(student.UUID)
		if err == nil {
			for _, skill := range skills {
				studentResponse.Skills = append(studentResponse.Skills, skill.Name)
			}
		}

		responses = append(responses, studentResponse)
	}

	// Calculate total pages
	totalPages := totalCount / req.Limit
	if totalCount%req.Limit != 0 {
		totalPages++
	}

	return &modelDTO.StudentListResponse{
		Students:   responses,
		TotalCount: totalCount,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}
