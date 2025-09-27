package delivery

import "github.com/HPNV/growlink-backend/delivery/developer"

type IDelivery interface {
	GetDeveloper() developer.IDeveloper
}

type Delivery struct {
	developer developer.IDeveloper
}

func NewDelivery(developerService developer.IDeveloper) IDelivery {
	return &Delivery{
		developer: developerService,
	}
}

func (d *Delivery) GetDeveloper() developer.IDeveloper {
	return d.developer
}
