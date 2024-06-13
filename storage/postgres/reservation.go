package postgres

import (
	"database/sql"
	"fmt"
	"reservation-service/client"
	"reservation-service/config/logger"
	pb "reservation-service/genproto/reservation"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReservationRepo struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewReservationRepo(db *sql.DB, logger *logger.Logger) *ReservationRepo {
	return &ReservationRepo{db: db, Logger: logger}
}

func (r *ReservationRepo) Create(req *pb.ReservationReq) (*pb.Reservation, error) {
	id := uuid.New().String()
	res := pb.Reservation{}

	query := `INSERT INTO reservations (
		id, 
		user_id, 
		restaurant_id, 
		reservation_time, 
		status
	) VALUES ($1, $2, $3, $4, $5) 
	 	RETURNING 
		id,
		user_id,
		restaurant_id,
		reservation_time,
		status`

	var reservationTime time.Time

	row := r.db.QueryRow(query, id, req.UserId, req.RestaurantId, req.ReservationTime, req.Status)
	err := row.Scan(
		&res.Id,
		&res.UserId,
		&res.RestaurantId,
		&reservationTime,
		&res.Status,
	)
	if err != nil {
		r.Logger.ERROR.Println("Error while creating reservation: ", err)
		return nil, err
	}

	req.ReservationTime = reservationTime.Format("2006-01-02")

	r.Logger.INFO.Println("Successfully created reservation")
	return &res, nil
}

func (r *ReservationRepo) Get(id *pb.GetByIdReq) (*pb.ReservationRes, error) {
	res := &pb.ReservationRes{
		User:       &pb.UserResp{},
		Restaurant: &pb.Restaurant{},
	}

	query := `SELECT
					r.id,
					rs.id as restaurant_id,
					rs.name,
					rs.address,
					rs.phone_number,
					rs.description,
					r.reservation_time,
					r.status
				FROM reservations r
				JOIN restaurants rs ON r.restaurant_id = rs.id
				WHERE r.id = $1 AND r.deleted_at=0`

	row := r.db.QueryRow(query, id.Id)

	var reservationTime time.Time

	err := row.Scan(
		&res.Id,
		&res.Restaurant.Id,
		&res.Restaurant.Name,
		&res.Restaurant.Address,
		&res.Restaurant.PhoneNumber,
		&res.Restaurant.Description,
		&reservationTime,
		&res.Status,
	)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting reservation by id : ", err)
		return nil, err
	}

	query = `SELECT user_id from reservations WHERE id = $1`
	row = r.db.QueryRow(query, id.Id)

	var userId string
	err = row.Scan(&userId)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting user id : ", err)
		return nil, err
	}

	us, err := client.GetUser(userId)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting user : ", err)
		return nil, err
	}

	res.User.Id = us.ID
	res.User.Username = us.Username
	res.User.Email = us.Email

	res.ReservationTime = reservationTime.Format("2006-01-02")

	return res, nil
}

func (r *ReservationRepo) GetAll(req *pb.GetAllReservationReq) (*pb.GetAllReservationRes, error) {
	res := &pb.GetAllReservationRes{
		Reservation: []*pb.ReservationRes{},
	}

	query := `SELECT
					r.id,
					rs.id as restaurant_id,
					rs.name,
					rs.address,
					rs.phone_number,
					rs.description,
					r.reservation_time,
					r.status
				FROM reservations r
				JOIN restaurants rs ON r.restaurant_id = rs.id
				WHERE r.deleted_at=0`

	var args []interface{}
	var conditions []string

	if req.UserId != "" {
		args = append(args, req.UserId)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var defaultLimit int32
	err := r.db.QueryRow("SELECT COUNT(1) FROM reservations WHERE deleted_at=0").Scan(&defaultLimit)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting count : ", err)
		return nil, err
	}
	if req.Filter.Limit == 0 {
		req.Filter.Limit = defaultLimit
	}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting all reservations : ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rs := &pb.ReservationRes{
			User:       &pb.UserResp{},
			Restaurant: &pb.Restaurant{},
		}

		var reservationTime time.Time

		err := rows.Scan(
			&rs.Id,
			&rs.Restaurant.Id,
			&rs.Restaurant.Name,
			&rs.Restaurant.Address,
			&rs.Restaurant.PhoneNumber,
			&rs.Restaurant.Description,
			&reservationTime,
			&rs.Status,
		)
		if err != nil {
			r.Logger.ERROR.Println("Error while scanning all reservations : ", err)
			return nil, err
		}

		query = `SELECT user_id from reservations WHERE id = $1`
		row := r.db.QueryRow(query, rs.Id)

		var userId string
		err = row.Scan(&userId)
		if err != nil {
			r.Logger.ERROR.Println("Error while getting user id : ", err)
			return nil, err
		}

		us, err := client.GetUser(userId)
		if err != nil {
			r.Logger.ERROR.Println("Error while getting user : ", err)
			return nil, err
		}

		rs.User.Id = us.ID
		rs.User.Username = us.Username
		rs.User.Email = us.Email

		rs.ReservationTime = reservationTime.Format("2006-01-02")

		res.Reservation = append(res.Reservation, rs)
	}

	r.Logger.INFO.Println("Successfully fetched all reservations")
	return res, nil
}

func (r *ReservationRepo) Update(req *pb.ReservationUpdate) (*pb.Reservation, error) {
	res := pb.Reservation{}

	query := `UPDATE reservations SET user_id=$1, restaurant_id=$2, reservation_time=$3, status=$4, updated_at=now() WHERE id=$5 RETURNING id, user_id, restaurant_id, reservation_time`

	row := r.db.QueryRow(query, req.UpdateReservation.UserId, req.UpdateReservation.RestaurantId, req.UpdateReservation.ReservationTime, req.UpdateReservation.Status, req.Id.Id)

	var reservationTime time.Time

	err := row.Scan(
		&res.Id,
		&res.UserId,
		&res.RestaurantId,
		&reservationTime,
	)
	if err != nil {
		r.Logger.ERROR.Println("Error while updating reservation : ", err)
		return nil, err
	}
	res.ReservationTime = reservationTime.Format("2006-01-02")

	r.Logger.INFO.Println("Successfully updated reservation")
	return &res, nil
}

func (r *ReservationRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := pb.Void{}

	query := `UPDATE reservations SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1`
	_, err := r.db.Exec(query, req.Id)
	if err != nil {
		r.Logger.ERROR.Println("Error while deleting reservation : ", err)
		return nil, err
	}

	r.Logger.INFO.Println("Successfully deleted reservation")
	return &res, nil
}
