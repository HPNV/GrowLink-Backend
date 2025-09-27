package repository

import (
	"github.com/HPNV/growlink-backend/repository/developer"
	"github.com/jmoiron/sqlx"
)

type IRegistry interface {
	GetDB() *sqlx.DB
	GetDeveloper() developer.IDeveloper
}

type Registry struct {
	db            *sqlx.DB
	developerRepo developer.IDeveloper
}

func NewRegistry(
	db *sqlx.DB,
	developerRepo developer.IDeveloper,
) *Registry {
	return &Registry{
		db:            db,
		developerRepo: developerRepo,
	}
}

func (r *Registry) GetDB() *sqlx.DB {
	return r.db
}

func (r *Registry) GetDeveloper() developer.IDeveloper {
	return r.developerRepo
}
