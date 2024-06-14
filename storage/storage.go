package storage

import (
	r "reservation-service/genproto/reservation"
)

type StorageI interface {
	Reservation() ReservationI
}

type RestaurantI interface {
	Create(*r.RestaurantReq) (*r.Restaurant, error)
	Get(*r.GetByIdReq) (*r.Restaurant, error)
	GetAll(*r.GetAllRestaurantReq) (*r.GetAllRestaurantRes, error)
	Update(*r.RestaurantUpdate) (*r.Restaurant, error)
	Delete(*r.GetByIdReq) (*r.Void, error)
}

type ReservationI interface {
	Create(*r.ReservationReq) (*r.Reservation, error)
	Get(*r.GetByIdReq) (*r.ReservationRes, error)
	GetAll(*r.GetAllReservationReq) (*r.GetAllReservationRes, error)
	Update(*r.ReservationUpdate) (*r.Reservation, error)
	Delete(*r.GetByIdReq) (*r.Void, error)
	CheckTime(req *r.CheckTimeReq) (*r.CheckTimeResp, error)
	GetMenu(req *r.GetMenuReq) (*r.GetAllMenuRess, error)
}

type ReservationOrderI interface {
	Create(*r.ReservationOrderReq) (*r.ReservationOrderRes, error)
	Get(*r.GetByIdReq) (*r.ReservationOrderRes, error)
	GetAll(*r.GetAllReservationOrderReq) (*r.GetAllReservationOrderRes, error)
	Update(*r.ReservationOrderUpdateReq) (*r.ReservationOrderRes, error)
	Delete(*r.GetByIdReq) (*r.Void, error)
}

type MenuI interface {
	Create(*r.MenuReq) (*r.Menu, error)
	Get(*r.GetByIdReq) (*r.MenuRes, error)
	GetAll(*r.GetAllMenuReq) (*r.GetAllMenuRes, error)
	Update(*r.MenuUpdate) (*r.Menu, error)
	Delete(*r.GetByIdReq) (*r.Void, error)
}
