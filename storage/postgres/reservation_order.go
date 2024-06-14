package postgres

import (
	"database/sql"
	"fmt"
	"reservation-service/client"
	"reservation-service/config/logger"
	r "reservation-service/genproto/reservation"
	"time"

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
	res := r.ReservationOrderRes{
		Reservation: &r.ReservationRes{
			User:       &r.UserResp{},
			Restaurant: &r.Restaurant{},
		},
		MenuItem: &r.MenuRes{
			Restaurant: &r.Restaurant{},
		},
	}

	tr, err := o.db.Begin()
	if err != nil {
		o.Logger.ERROR.Println("Error beginning transaction for company deletion", "id", id, "error", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tr.Rollback()
			o.Logger.ERROR.Println("Panic during transaction for company deletion", "id", id, "panic", p)
			panic(p)
		} else if err != nil {
			tr.Rollback()
			o.Logger.ERROR.Println("Error during transaction for company deletion", "id", id, "error", err)
		} else {
			err = tr.Commit()
			if err != nil {
				o.Logger.ERROR.Println("Error committing transaction for company deletion", "id", id, "error", err)
			} else {
				o.Logger.ERROR.Println("Successfully committed transaction for company deletion", "id", id)
			}
		}
	}()

	insertQuery := `INSERT INTO reservation_orders (
		id,
		reservation_id,
		menu_item_id,
		quantity
	) VALUES ($1, $2, $3, $4)`

	_, err = o.db.Exec(insertQuery, id, req.ReservationId, req.MenuItemId, req.Quantity)
	if err != nil {
		o.Logger.ERROR.Println("Error while inserting reservation_order:", err)
		return nil, err
	}

	selectQuery := `SELECT
				ro.id,
				r.id as reservation_id,
				rs.id as restaurant_id,
				rs.name as restaurant_name,
				rs.address,
				rs.phone_number,
				rs.description,
				r.reservation_time,
				r.status,
				m.id as menu_item_id,
				m.name as menu_item_name,
				m.description as menu_item_description,
				m.price as menu_item_price,
				rsm.id as restaurant_id,
				rsm.name as restaurant_name,
				rsm.address,
				rsm.phone_number,
				rsm.description,
				ro.quantity
			FROM reservation_orders ro 
			JOIN reservations r ON ro.reservation_id = r.id
			JOIN restaurants rs ON r.restaurant_id = rs.id
			JOIN menu m ON ro.menu_item_id = m.id
			JOIN restaurants rsm ON m.restaurant_id = rsm.id
			WHERE ro.id = $1 AND ro.deleted_at = 0`

	var reservationTime time.Time

	row := tr.QueryRow(selectQuery, id)
	err = row.Scan(
		&res.Id,
		&res.Reservation.Id,
		&res.Reservation.Restaurant.Id,
		&res.Reservation.Restaurant.Name,
		&res.Reservation.Restaurant.Address,
		&res.Reservation.Restaurant.PhoneNumber,
		&res.Reservation.Restaurant.Description,
		&reservationTime,
		&res.Reservation.Status,
		&res.MenuItem.Id,
		&res.MenuItem.Name,
		&res.MenuItem.Description,
		&res.MenuItem.Price,
		&res.MenuItem.Restaurant.Id,
		&res.MenuItem.Restaurant.Name,
		&res.MenuItem.Restaurant.Address,
		&res.MenuItem.Restaurant.PhoneNumber,
		&res.MenuItem.Restaurant.Description,
		&res.Quantity,
	)
	if err != nil {
		o.Logger.ERROR.Println("Error while retrieving reservation_order details:", err)
		return nil, err
	}

	userQuery := `SELECT user_id FROM reservations WHERE id = $1`
	row = tr.QueryRow(userQuery, res.Reservation.Id)

	var userId string
	err = row.Scan(&userId)
	if err != nil {
		o.Logger.ERROR.Println("Error while getting user id:", err)
		return nil, err
	}

	us, err := client.GetUser(userId)
	if err != nil {
		o.Logger.ERROR.Println("Error while getting user:", err)
		return nil, err
	}

	res.Reservation.User.Id = us.ID
	res.Reservation.User.Username = us.Username
	res.Reservation.User.Email = us.Email
	res.Reservation.ReservationTime = reservationTime.Format("2006-01-02 15:04:05")

	if res.Reservation.Restaurant.Id != res.MenuItem.Restaurant.Id {
		o.Logger.ERROR.Println("Restaurants are not the same!!!")
		tr.Rollback()
		return nil, fmt.Errorf("Restaurants are not the same!!!")
	}

	o.Logger.INFO.Println("Successfully created reservation_order")

	fmt.Println(res)
	return &res, nil
}

func (o *ReservationOrderRepo) Get(req *r.GetByIdReq) (*r.ReservationOrderRes, error) {
	res := &r.ReservationOrderRes{
		Reservation: &r.ReservationRes{
			User:       &r.UserResp{},
			Restaurant: &r.Restaurant{},
		},
		MenuItem: &r.MenuRes{
			Restaurant: &r.Restaurant{},
		},
	}

	selectQuery := `SELECT
				ro.id,
				r.id as reservation_id,
				rs.id as restaurant_id,
				rs.name as restaurant_name,
				rs.address,
				rs.phone_number,
				rs.description,
				r.reservation_time,
				r.status,
				m.id as menu_item_id,
				m.name as menu_item_name,
				m.description as menu_item_description,
				m.price as menu_item_price,
				rsm.id as restaurant_id,
				rsm.name as restaurant_name,
				rsm.address,
				rsm.phone_number,
				rsm.description,
				ro.quantity
			FROM reservation_orders ro 
			JOIN reservations r ON ro.reservation_id = r.id
			JOIN restaurants rs ON r.restaurant_id = rs.id
			JOIN menu m ON ro.menu_item_id = m.id
			JOIN restaurants rsm ON m.restaurant_id = rsm.id
			WHERE ro.id = $1 AND ro.deleted_at = 0`

	var reservationTime time.Time

	row := o.db.QueryRow(selectQuery, req.Id)
	err := row.Scan(
		&res.Id,
		&res.Reservation.Id,
		&res.Reservation.Restaurant.Id,
		&res.Reservation.Restaurant.Name,
		&res.Reservation.Restaurant.Address,
		&res.Reservation.Restaurant.PhoneNumber,
		&res.Reservation.Restaurant.Description,
		&reservationTime,
		&res.Reservation.Status,
		&res.MenuItem.Id,
		&res.MenuItem.Name,
		&res.MenuItem.Description,
		&res.MenuItem.Price,
		&res.MenuItem.Restaurant.Id,
		&res.MenuItem.Restaurant.Name,
		&res.MenuItem.Restaurant.Address,
		&res.MenuItem.Restaurant.PhoneNumber,
		&res.MenuItem.Restaurant.Description,
		&res.Quantity,
	)
	if err != nil {
		o.Logger.ERROR.Println("Error while retrieving reservation_order details:", err)
		return nil, err
	}

	userQuery := `SELECT user_id FROM reservations WHERE id = $1`
	row = o.db.QueryRow(userQuery, res.Reservation.Id)

	var userId string
	err = row.Scan(&userId)
	if err != nil {
		o.Logger.ERROR.Println("Error while getting user id:", err)
		return nil, err
	}

	us, err := client.GetUser(userId)
	if err != nil {
		o.Logger.ERROR.Println("Error while getting user:", err)
		return nil, err
	}

	res.Reservation.User.Id = us.ID
	res.Reservation.User.Username = us.Username
	res.Reservation.User.Email = us.Email
	res.Reservation.ReservationTime = reservationTime.Format("2006-01-02 15:04:05")

	return res, nil
}

func (o *ReservationOrderRepo) GetAll(req *r.GetAllReservationOrderReq) (*r.GetAllReservationOrderRes, error) {
	res := &r.GetAllReservationOrderRes{
		ReservationOrder: []*r.ReservationOrderRes{},
	}

	selectQuery := `SELECT
				ro.id,
				r.id as reservation_id,
				rs.id as restaurant_id,
				rs.name as restaurant_name,
				rs.address,
				rs.phone_number,
				rs.description,
				r.reservation_time,
				r.status,
				m.id as menu_item_id,
				m.name as menu_item_name,
				m.description as menu_item_description,
				m.price as menu_item_price,
				rsm.id as restaurant_id,
				rsm.name as restaurant_name,
				rsm.address,
				rsm.phone_number,
				rsm.description,
				ro.quantity
			FROM reservation_orders ro 
			JOIN reservations r ON ro.reservation_id = r.id
			JOIN restaurants rs ON r.restaurant_id = rs.id
			JOIN menu m ON ro.menu_item_id = m.id
			JOIN restaurants rsm ON m.restaurant_id = rsm.id
			WHERE ro.deleted_at = 0 LIMIT $1 OFFSET $2`

	rows, err := o.db.Query(selectQuery, req.Filter.Limit, req.Filter.Offset)
	if err != nil {
		o.Logger.ERROR.Println("Error while retrieving reservation_orders:", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		resOrd := &r.ReservationOrderRes{
			Reservation: &r.ReservationRes{
				User:       &r.UserResp{},
				Restaurant: &r.Restaurant{},
			},
			MenuItem: &r.MenuRes{
				Restaurant: &r.Restaurant{},
			},
		}
		var reservationTime time.Time
		err := rows.Scan(
			&resOrd.Id,
			&resOrd.Reservation.Id,
			&resOrd.Reservation.Restaurant.Id,
			&resOrd.Reservation.Restaurant.Name,
			&resOrd.Reservation.Restaurant.Address,
			&resOrd.Reservation.Restaurant.PhoneNumber,
			&resOrd.Reservation.Restaurant.Description,
			&reservationTime,
			&resOrd.Reservation.Status,
			&resOrd.MenuItem.Id,
			&resOrd.MenuItem.Name,
			&resOrd.MenuItem.Description,
			&resOrd.MenuItem.Price,
			&resOrd.MenuItem.Restaurant.Id,
			&resOrd.MenuItem.Restaurant.Name,
			&resOrd.MenuItem.Restaurant.Address,
			&resOrd.MenuItem.Restaurant.PhoneNumber,
			&resOrd.MenuItem.Restaurant.Description,
			&resOrd.Quantity,
		)
		if err != nil {
			o.Logger.ERROR.Println("Error while retrieving reservation_order details:", err)
			return nil, err
		}

		userQuery := `SELECT user_id FROM reservations WHERE id = $1`
		row := o.db.QueryRow(userQuery, resOrd.Reservation.Id)

		var userId string
		err = row.Scan(&userId)
		if err != nil {
			o.Logger.ERROR.Println("Error while getting user id:", err)
			return nil, err
		}

		us, err := client.GetUser(userId)
		if err != nil {
			o.Logger.ERROR.Println("Error while getting user:", err)
			return nil, err
		}

		resOrd.Reservation.User.Id = us.ID
		resOrd.Reservation.User.Username = us.Username
		resOrd.Reservation.User.Email = us.Email
		resOrd.Reservation.ReservationTime = reservationTime.Format("2006-01-02 15:04:05")

		res.ReservationOrder = append(res.ReservationOrder, resOrd)
	}

	return res, nil
}

func (o *ReservationOrderRepo) Update(req *r.ReservationOrderUpdateReq) (*r.ReservationOrderRes, error) {
	res := r.ReservationOrderRes{}

	query := `UPDATE reservation_orders SET reservation_id=$1, menu_item_id=$2, quantity=$3 WHERE id=$4 and deleted_at=0`

	_, err := o.db.Exec(query, req.Update.ReservationId, req.Update.MenuItemId, req.Update.Quantity, req.Id.Id)
	if err != nil {
		o.Logger.ERROR.Println("Error while updating reservation_order:", err)
		return nil, err
	}

	selectQuery := `SELECT
				ro.id,
				r.id as reservation_id,
				rs.id as restaurant_id,
				rs.name as restaurant_name,
				rs.address,
				rs.phone_number,
				rs.description,
				r.reservation_time,
				r.status,
				m.id as menu_item_id,
				m.name as menu_item_name,
				m.description as menu_item_description,
				m.price as menu_item_price,
				rsm.id as restaurant_id,
				rsm.name as restaurant_name,
				rsm.address,
				rsm.phone_number,
				rsm.description,
				ro.quantity
			FROM reservation_orders ro 
			JOIN reservations r ON ro.reservation_id = r.id
			JOIN restaurants rs ON r.restaurant_id = rs.id
			JOIN menu m ON ro.menu_item_id = m.id
			JOIN restaurants rsm ON m.restaurant_id = rsm.id
			WHERE ro.id = $1 AND ro.deleted_at = 0`

	var reservationTime time.Time

	row := o.db.QueryRow(selectQuery, req.Id)
	err = row.Scan(
		&res.Id,
		&res.Reservation.Id,
		&res.Reservation.Restaurant.Id,
		&res.Reservation.Restaurant.Name,
		&res.Reservation.Restaurant.Address,
		&res.Reservation.Restaurant.PhoneNumber,
		&res.Reservation.Restaurant.Description,
		&reservationTime,
		&res.Reservation.Status,
		&res.MenuItem.Id,
		&res.MenuItem.Name,
		&res.MenuItem.Description,
		&res.MenuItem.Price,
		&res.MenuItem.Restaurant.Id,
		&res.MenuItem.Restaurant.Name,
		&res.MenuItem.Restaurant.Address,
		&res.MenuItem.Restaurant.PhoneNumber,
		&res.MenuItem.Restaurant.Description,
		&res.Quantity,
	)
	if err != nil {
		o.Logger.ERROR.Println("Error while retrieving reservation_order details:", err)
		return nil, err
	}

	userQuery := `SELECT user_id FROM reservations WHERE id = $1`
	row = o.db.QueryRow(userQuery, res.Reservation.Id)

	var userId string
	err = row.Scan(&userId)
	if err != nil {
		o.Logger.ERROR.Println("Error while getting user id:", err)
		return nil, err
	}

	us, err := client.GetUser(userId)
	if err != nil {
		o.Logger.ERROR.Println("Error while getting user:", err)
		return nil, err
	}

	res.Reservation.User.Id = us.ID
	res.Reservation.User.Username = us.Username
	res.Reservation.User.Email = us.Email
	res.Reservation.ReservationTime = reservationTime.Format("2006-01-02 15:04:05")

	return &res, nil
}

func (o *ReservationOrderRepo) Delete(req *r.GetByIdReq) (*r.Void, error) {
	res := r.Void{}

	query := `UPDATE reservation_orders SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1`
	ress, err := o.db.Exec(query, req.Id)
	if err != nil {
		o.Logger.ERROR.Println("Error while deleting reservation_order")
		return nil, err
	}

	if r, err := ress.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("restaurant with id %s not found", req.Id)
	}

	return &res, nil
}
