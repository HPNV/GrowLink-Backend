package developer

import (
	modelDB "github.com/HPNV/growlink-backend/model/db"
	"github.com/jmoiron/sqlx"
)

type IDeveloper interface {
	CreateDeveloper(developer *modelDB.Developer) error
}

type Developer struct {
	DB *sqlx.DB
}

func NewDeveloperRepo(db *sqlx.DB) *Developer {
	return &Developer{
		DB: db,
	}
}

func (r *Developer) CreateDeveloper(developer *modelDB.Developer) error {
	_, err := r.DB.Exec(InsertDeveloperQuery, developer.UUID, developer.Name, developer.Email)
	if err != nil {
		return err
	}

	return nil
}
