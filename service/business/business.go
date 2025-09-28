package business

import (
	"github.com/HPNV/growlink-backend/model/db"
	"github.com/HPNV/growlink-backend/model/dto"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/jmoiron/sqlx"
)

type IBusiness interface {
	Create(userUUID string, req *dto.BusinessRequest) (*dto.BusinessResponse, error)
	GetByUUID(uuid string) (*dto.BusinessResponse, error)
	GetByUserUUID(userUUID string) (*dto.BusinessResponse, error)
	Update(uuid string, req *dto.BusinessRequest) (*dto.BusinessResponse, error)
	Delete(uuid string) error
	GetAll() ([]*dto.BusinessResponse, error)
}

type Business struct {
	repo repository.IRegistry
}

func NewBusiness(repo repository.IRegistry) IBusiness {
	return &Business{
		repo: repo,
	}
}

func (b *Business) Create(userUUID string, req *dto.BusinessRequest) (*dto.BusinessResponse, error) {
	business := &db.Business{
		UserUUID:    userUUID,
		CompanyName: req.CompanyName,
	}

	err := b.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return b.repo.GetBusiness().Create(tx, business)
	})

	if err != nil {
		return nil, err
	}

	return &dto.BusinessResponse{
		UUID:        business.UUID,
		UserUUID:    business.UserUUID,
		CompanyName: business.CompanyName,
	}, nil
}

func (b *Business) GetByUUID(uuid string) (*dto.BusinessResponse, error) {
	business, err := b.repo.GetBusiness().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	return &dto.BusinessResponse{
		UUID:        business.UUID,
		UserUUID:    business.UserUUID,
		CompanyName: business.CompanyName,
	}, nil
}

func (b *Business) GetByUserUUID(userUUID string) (*dto.BusinessResponse, error) {
	business, err := b.repo.GetBusiness().GetByUserUUID(userUUID)
	if err != nil {
		return nil, err
	}

	return &dto.BusinessResponse{
		UUID:        business.UUID,
		UserUUID:    business.UserUUID,
		CompanyName: business.CompanyName,
	}, nil
}

func (b *Business) Update(uuid string, req *dto.BusinessRequest) (*dto.BusinessResponse, error) {
	// First get existing business
	existing, err := b.repo.GetBusiness().GetByUUID(uuid)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.CompanyName = req.CompanyName

	err = b.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return b.repo.GetBusiness().Update(tx, existing)
	})

	if err != nil {
		return nil, err
	}

	return &dto.BusinessResponse{
		UUID:        existing.UUID,
		UserUUID:    existing.UserUUID,
		CompanyName: existing.CompanyName,
	}, nil
}

func (b *Business) Delete(uuid string) error {
	return b.repo.WithTransaction(func(tx *sqlx.Tx) error {
		return b.repo.GetBusiness().Delete(tx, uuid)
	})
}

func (b *Business) GetAll() ([]*dto.BusinessResponse, error) {
	businesses, err := b.repo.GetBusiness().GetAll()
	if err != nil {
		return nil, err
	}

	var responses []*dto.BusinessResponse
	for _, business := range businesses {
		responses = append(responses, &dto.BusinessResponse{
			UUID:        business.UUID,
			UserUUID:    business.UserUUID,
			CompanyName: business.CompanyName,
		})
	}

	return responses, nil
}
