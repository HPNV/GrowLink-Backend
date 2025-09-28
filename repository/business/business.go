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
	query := `
		INSERT INTO businesses (user_uuid, company_name)
		VALUES ($1, $2)
		RETURNING uuid
	`
	return tx.QueryRow(query, business.UserUUID, business.CompanyName).Scan(&business.UUID)
}

func (b *Business) GetByUUID(uuid string) (*db.Business, error) {
	business := &db.Business{}
	query := `SELECT uuid, user_uuid, company_name FROM businesses WHERE uuid = $1`
	err := b.db.Get(business, query, uuid)
	return business, err
}

func (b *Business) GetByUserUUID(userUUID string) (*db.Business, error) {
	business := &db.Business{}
	query := `SELECT uuid, user_uuid, company_name FROM businesses WHERE user_uuid = $1`
	err := b.db.Get(business, query, userUUID)
	return business, err
}

func (b *Business) Update(tx *sqlx.Tx, business *db.Business) error {
	query := `
		UPDATE businesses 
		SET company_name = $1
		WHERE uuid = $2
	`
	_, err := tx.Exec(query, business.CompanyName, business.UUID)
	return err
}

func (b *Business) Delete(tx *sqlx.Tx, uuid string) error {
	query := `DELETE FROM businesses WHERE uuid = $1`
	_, err := tx.Exec(query, uuid)
	return err
}

func (b *Business) GetAll() ([]*db.Business, error) {
	var businesses []*db.Business
	query := `SELECT uuid, user_uuid, company_name FROM businesses ORDER BY company_name`
	err := b.db.Select(&businesses, query)
	return businesses, err
}
