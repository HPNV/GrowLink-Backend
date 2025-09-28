package service

import (
	"github.com/HPNV/growlink-backend/service/business"
	"github.com/HPNV/growlink-backend/service/project"
	"github.com/HPNV/growlink-backend/service/skill"
	"github.com/HPNV/growlink-backend/service/student"
	"github.com/HPNV/growlink-backend/service/user"
)

type IRegistry interface {
	GetUser() user.IUser
	GetSkill() skill.ISkill
	GetBusiness() business.IBusiness
	GetStudent() student.IStudent
	GetProject() project.IProject
}

type Registry struct {
	user     user.IUser
	skill    skill.ISkill
	business business.IBusiness
	student  student.IStudent
	project  project.IProject
}

func NewRegistry(
	user user.IUser,
	skill skill.ISkill,
	business business.IBusiness,
	student student.IStudent,
	project project.IProject,
) *Registry {
	return &Registry{
		user:     user,
		skill:    skill,
		business: business,
		student:  student,
		project:  project,
	}
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
