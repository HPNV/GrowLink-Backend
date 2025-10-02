package business

import (
	"github.com/HPNV/growlink-backend/model/db"
	"github.com/jmoiron/sqlx"
)

type IBusiness interface {
	Create(tx *sqlx.Tx, business *db.Business) error
	GetByUUID(uuid string) (*db.Business, error)
	GetByUserUUID(userUUID string) (*db.Business, error)
	Update(tx *sqlx.Tx, business *db.Business) error
	Delete(tx *sqlx.Tx, uuid string) error
	GetAll() ([]*db.Business, error)
}

type Business struct {
	db *sqlx.DB
}

func NewBusiness(db *sqlx.DB) IBusiness {
	return &Business{
		db: db,
	}
}

func (b *Business) Create(tx *sqlx.Tx, business *db.Business) error {
	return tx.QueryRow(CreateQuery, business.UserUUID, business.CompanyName).Scan(&business.UUID)
}

func (b *Business) GetByUUID(uuid string) (*db.Business, error) {
	business := &db.Business{}
	err := b.db.Get(business, GetByUUIDQuery, uuid)
	return business, err
}

func (b *Business) GetByUserUUID(userUUID string) (*db.Business, error) {
	business := &db.Business{}
	err := b.db.Get(business, GetByUserUUIDQuery, userUUID)
	return business, err
}

func (b *Business) Update(tx *sqlx.Tx, business *db.Business) error {
	_, err := tx.Exec(UpdateQuery, business.CompanyName, business.UUID)
	return err
}

func (b *Business) Delete(tx *sqlx.Tx, uuid string) error {
	_, err := tx.Exec(DeleteQuery, uuid)
	return err
}

func (b *Business) GetAll() ([]*db.Business, error) {
	var businesses []*db.Business
	err := b.db.Select(&businesses, GetAllQuery)
	return businesses, err
}
