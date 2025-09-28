package delivery

import (
	"github.com/HPNV/growlink-backend/delivery/business"
	"github.com/HPNV/growlink-backend/delivery/file"
	"github.com/HPNV/growlink-backend/delivery/project"
	"github.com/HPNV/growlink-backend/delivery/skill"
	"github.com/HPNV/growlink-backend/delivery/student"
	"github.com/HPNV/growlink-backend/delivery/user"
)

type IDelivery interface {
	GetUser() user.IUser
	GetBusiness() business.IBusiness
	GetStudent() student.IStudent
	GetSkill() skill.ISkill
	GetProject() project.IProject
	GetFile() file.IFile
}

type Delivery struct {
	user     user.IUser
	business business.IBusiness
	student  student.IStudent
	skill    skill.ISkill
	project  project.IProject
	file     file.IFile
}

func NewDelivery(
	user user.IUser,
	business business.IBusiness,
	student student.IStudent,
	skill skill.ISkill,
	project project.IProject,
	file file.IFile,
) IDelivery {
	return &Delivery{
		user:     user,
		business: business,
		student:  student,
		skill:    skill,
		project:  project,
		file:     file,
	}
}

func (d *Delivery) GetUser() user.IUser {
	return d.user
}

func (d *Delivery) GetBusiness() business.IBusiness {
	return d.business
}

func (d *Delivery) GetStudent() student.IStudent {
	return d.student
}

func (d *Delivery) GetSkill() skill.ISkill {
	return d.skill
}

func (d *Delivery) GetProject() project.IProject {
	return d.project
}

func (d *Delivery) GetFile() file.IFile {
	return d.file
}
