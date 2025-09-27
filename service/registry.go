package service

import "github.com/HPNV/growlink-backend/service/developer"

type IRegistry interface {
	GetDeveloper() developer.IDeveloper
}

type Registry struct {
	developer developer.IDeveloper
}

func NewRegistry(developer developer.IDeveloper) *Registry {
	return &Registry{
		developer: developer,
	}
}

func (r *Registry) GetDeveloper() developer.IDeveloper {
	return r.developer
}
