package delivery

import (
	"github.com/HPNV/growlink-backend/delivery/business"
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
}

type Delivery struct {
	user     user.IUser
	business business.IBusiness
	student  student.IStudent
	skill    skill.ISkill
	project  project.IProject
}

func NewDelivery(
	user user.IUser,
	business business.IBusiness,
	student student.IStudent,
	skill skill.ISkill,
	project project.IProject,
) IDelivery {
	return &Delivery{
		user:     user,
		business: business,
		student:  student,
		skill:    skill,
		project:  project,
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
