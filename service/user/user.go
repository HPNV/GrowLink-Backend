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
