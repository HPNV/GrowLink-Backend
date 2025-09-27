package delivery

import "github.com/HPNV/growlink-backend/delivery/user"

type IDelivery interface {
	GetUser() user.IUser
}

type Delivery struct {
	user user.IUser
}

func NewDelivery(
	user user.IUser,
) IDelivery {
	return &Delivery{
		user: user,
	}
}

func (d *Delivery) GetUser() user.IUser {
	return d.user
}
