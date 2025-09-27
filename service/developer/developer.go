package developer

import (
	modelDB "github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/repository"
)

type IDeveloper interface {
	CreateDeveloper(developer *modelDB.Developer) error
}

type Developer struct {
	repo repository.IRegistry
}

func NewDeveloperService(repo repository.IRegistry) *Developer {
	return &Developer{
		repo: repo,
	}
}

func (s *Developer) CreateDeveloper(developer *modelDB.Developer) error {
	err := s.repo.GetDeveloper().CreateDeveloper(developer)
	if err != nil {
		return err
	}

	return nil
}
