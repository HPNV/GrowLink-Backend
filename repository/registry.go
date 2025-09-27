package repository

import (
	"github.com/HPNV/growlink-backend/repository/user"
	"github.com/jmoiron/sqlx"
)

type IRegistry interface {
	GetDB() *sqlx.DB
	GetUser() user.IUser
}

type Registry struct {
	db   *sqlx.DB
	user user.IUser
}

func NewRegistry(
	db *sqlx.DB,
	user user.IUser,
) *Registry {
	return &Registry{
		db:   db,
		user: user,
	}
}

func (r *Registry) GetDB() *sqlx.DB {
	return r.db
}

func (r *Registry) GetUser() user.IUser {
	return r.user
}
