package repository

import (
	"github.com/HPNV/growlink-backend/repository/business"
	"github.com/HPNV/growlink-backend/repository/file"
	"github.com/HPNV/growlink-backend/repository/project"
	"github.com/HPNV/growlink-backend/repository/skill"
	"github.com/HPNV/growlink-backend/repository/student"
	"github.com/HPNV/growlink-backend/repository/user"
	"github.com/jmoiron/sqlx"
)

type IRegistry interface {
	GetDB() *sqlx.DB
	GetUser() user.IUser
	GetSkill() skill.ISkill
	GetBusiness() business.IBusiness
	GetStudent() student.IStudent
	GetProject() project.IProject
	GetFile() file.IFile
	WithTransaction(txFunc func(tx *sqlx.Tx) error) error
}

type Registry struct {
	db       *sqlx.DB
	user     user.IUser
	skill    skill.ISkill
	business business.IBusiness
	student  student.IStudent
	project  project.IProject
	file     file.IFile
}

func NewRegistry(
	db *sqlx.DB,
	user user.IUser,
	skill skill.ISkill,
	business business.IBusiness,
	student student.IStudent,
	project project.IProject,
	file file.IFile,
) *Registry {
	return &Registry{
		db:       db,
		user:     user,
		skill:    skill,
		business: business,
		student:  student,
		project:  project,
		file:     file,
	}
}

func (r *Registry) GetDB() *sqlx.DB {
	return r.db
}

func (r *Registry) GetUser() user.IUser {
	return r.user
}

func (r *Registry) GetSkill() skill.ISkill {
	return r.skill
}

func (r *Registry) GetBusiness() business.IBusiness {
	return r.business
}

func (r *Registry) GetStudent() student.IStudent {
	return r.student
}

func (r *Registry) GetProject() project.IProject {
	return r.project
}

func (r *Registry) GetFile() file.IFile {
	return r.file
}

func (r *Registry) WithTransaction(txFunc func(tx *sqlx.Tx) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()
	if err := txFunc(tx); err != nil {
		return err
	}
	return tx.Commit()
}
