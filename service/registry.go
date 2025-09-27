package service

import "github.com/HPNV/growlink-backend/service/user"

type IRegistry interface {
	GetUser() user.IUser
}

type Registry struct {
	user user.IUser
}

func NewRegistry(
	user user.IUser,
) *Registry {
	return &Registry{
		user: user,
	}
}

func (r *Registry) GetUser() user.IUser {
	return r.user
}
