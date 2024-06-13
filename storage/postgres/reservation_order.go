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

func (o *ReservationOrderRepo) Create(req *r.ReservationOrderReq) (*r.ReservationOrderRes, error) {
	id := uuid.New().String()
	res := r.ReservationOrderRes{}

	query := `INSERT INTO reservations (
		id,
		reservation_id,
		menu_item_id,
		quantity
	) VALUES ($1, $2, $3, $4) 
	RETURNING 
		id,
		reservation_id,
		menu_item_id,
		quantity`

	row := o.db.QueryRow(query, id, req.ReservationId, req.MenuItemId, req.Quantity)
	err := row.Scan(
		&res.Id,
		&res.Reservation.Id,
		&res.MenuItem.Id,
		&res.Quantity,
	)
	if err != nil {
		o.Logger.ERROR.Println("Error while creating reservation_orders")
		return nil, err
	}

	o.Logger.INFO.Println("Successfully created reservation_orders")

	return &res, nil
}

func (o *ReservationOrderRepo) Get(*r.GetByIdReq) (*r.ReservationOrderRes, error) {
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	return nil, nil
}

func (o *ReservationOrderRepo) GetAll(*r.GetAllReservationOrderReq) (*r.GetAllReservationOrderRes, error) {
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	// |---------------------------------------------------------------|
	return nil, nil
}

func(o *ReservationOrderRepo) Update(req *r.ReservationOrderUpdate) (*r.ReservationOrderRes, error){
	res := r.ReservationOrderRes{}

	query := `UPDATE reservation_orders SET reservation_id=$1, menu_item_id=$2, quantity=$3 WHERE id=$4 and deleted_at=0 RETURNING id, reservation_id, menu_item_id, quantity`

	row := o.db.QueryRow(query, req.UpdateBody.ReservationId, req.UpdateBody.MenuItemId, req.UpdateBody.Quantity, req.Id.Id)

	err := row.Scan(
		&res.Id,
		&res.Reservation.Id,
		&res.MenuItem.Id,
		&res.Quantity,
	)
	if err != nil {
		o.Logger.ERROR.Println("Error while updating reservation")
		return nil, err
	}

	return &res, nil
}

func (o *ReservationOrderRepo) Delete(req *r.GetByIdReq) (*r.Void, error) {
	res := r.Void{}

	query := `UPDATE reservation_orders SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1`
	_, err := o.db.Exec(query, req.Id)
	if err != nil {
		o.Logger.ERROR.Println("Error while deleting reservation_order")
		return nil, err
	}

	return &res, nil
}
