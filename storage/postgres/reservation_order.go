package postgres

import (
	"database/sql"
	"reservation-service/config/logger"
	r "reservation-service/genproto/reservation"

	"github.com/google/uuid"
)


type ReservationOrderRepo struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewReservationOrderRepo(db *sql.DB, logger *logger.Logger) *ReservationOrderRepo {
	return &ReservationOrderRepo{db: db, Logger: logger}
}

func(o *ReservationOrderRepo) Create(*r.ReservationOrderReq) (*r.ReservationOrderRes, error){
	return nil, nil
}

func(o *ReservationOrderRepo) Get(*r.GetByIdReq) (*r.ReservationOrderRes, error){
	return nil, nil
}

func(o *ReservationOrderRepo) GetAll(*r.GetAllReservationOrderReq) (*r.GetAllReservationOrderRes, error){
	return nil, nil
}

func(o *ReservationOrderRepo) Update(*r.ReservationOrderUpdate) (*r.ReservationOrderRes, error){
	return nil, nil
}

func(o *ReservationOrderRepo) Delete(*r.GetByIdReq) (*r.Void, error){
	return nil, nil
}
