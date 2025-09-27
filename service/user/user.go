package user

import (
	"context"

	modelDB "github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/repository"
)

type IUser interface {
	Login(ctx context.Context, email, password string) (*modelDB.User, error)
	Register(ctx context.Context, email, password, role string) (*modelDB.User, error)
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

func (u *User) Register(ctx context.Context, email, password, role string) (*modelDB.User, error) {
	user := &modelDB.User{
		Email: email,
		Role:  role,
	}
	user, err := u.repo.GetUser().Register(ctx, user, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
