package postgres

import (
	"database/sql"
	"fmt"
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

func NewReservationRepo(db *sql.DB) *ReservationRepo {
	return &ReservationRepo{db: db}
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
	 	RETURNING *`

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
		r.Logger.ERROR.Println("Error while creating reservation")
		return nil, err
	}

	req.ReservationTime = reservationTime.Format("2006-01-02")

	r.Logger.INFO.Panicln("Successfully created reservation")
	return &res, nil
}

func (r *ReservationRepo) Get(id *pb.GetByIdReq) (*pb.ReservationRes, error) {
	res := pb.ReservationRes{}
	query := `SELECT 
					r.id, 
					u.id, 
					u.username, 
					u.email, 
					rs.id, 
					rs.name, 
					rs.address, 
					rs.phone_number, 
					rs.description, 
					r.reservation_time, 
					r.status 
				FROM reservations r
				JOIN users u ON r.user_id = u.id 
				JOIN restaurants rs ON r.restaurant_id = rs.id
				WHERE r.id = $1 AND deleted_at=0`

	row := r.db.QueryRow(query, id.Id)

	var reservationTime time.Time

	err := row.Scan(
		&res.Id,
		&res.User.Id,
		&res.User.Username,
		&res.User.Email,
		&res.Restaurant.Id,
		&res.Restaurant.Name,
		&res.Restaurant.Address,
		&res.Restaurant.PhoneNumber,
		&res.Restaurant.Description,
		&reservationTime,
		&res.Status,
	)

	if err != nil {
		r.Logger.ERROR.Println("Error while getting reservation by id")
		return nil, err
	}

	res.ReservationTime = reservationTime.Format("2006-01-02")

	return &res, nil
}

func (r *ReservationRepo) GetAll(req *pb.GetAllReservationReq) (*pb.GetAllReservationRes, error) {
	res := pb.GetAllReservationRes{
		Reservation: []*pb.ReservationRes{},
	}

    query := `SELECT 
                    r.id, 
                    u.id, 
                    u.username, 
                    u.email, 
                    rs.id, 
                    rs.name, 
                    rs.address, 
                    rs.phone_number, 
                    rs.description, 
                    r.reservation_time, 
                    r.status 
                FROM reservations r
                JOIN users u ON r.user_id = u.id 
                JOIN restaurants rs ON r.restaurant_id = rs.id`

	var args []interface{}
	var conditions []string

	if req.UserId != "" {
		args = append(args, req.UserId)
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var defaultLimit int32
	err := r.db.QueryRow("SELECT COUNT(1) FROM reservation WHERE deleted_at=0").Scan(&defaultLimit)
	if err != nil {
		r.Logger.ERROR.Println("Error while get count")
		return nil, err
	}
	if req.Filter.Limit == 0 {
		req.Filter.Limit = defaultLimit
	}

	fmt.Println(query, args)

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting all reservations")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rs := pb.ReservationRes{}

		var reservationTime time.Time

		err := rows.Scan(
            &rs.Id,
			&rs.User.Id,
			&rs.User.Username,
			&rs.User.Email,
			&rs.Restaurant.Id,
			&rs.Restaurant.Name,
			&rs.Restaurant.Address,
			&rs.Restaurant.PhoneNumber,
			&rs.Restaurant.Description,
			&reservationTime,
			&rs.Status,
		)		
		if err != nil {
			r.Logger.ERROR.Println("Error while scan all reservations")
			return nil, err
		}

		rs.ReservationTime = reservationTime.Format("2006-01-02")

		res.Reservation = append(res.Reservation, &rs)
	}

	return &res, nil
}

func (r *ReservationRepo) Update (req *pb.ReservationUpdate) (*pb.Reservation, error) {
	res := pb.Reservation{}

	query := `UPDATE reservation SET user_id=$1, restaurant_id=$2, reservation_time=$3, status=$4 WHERE id=$5 RETURNING *`

	row := r.db.QueryRow(query, req.UpdateReservation.UserId, req.UpdateReservation.RestaurantId, req.UpdateReservation.ReservationTime, req.UpdateReservation.Status, req.Id)
	
	var reservationTime time.Time

	err := row.Scan(
        &res.Id,
        &res.UserId,
        &res.RestaurantId,
        &reservationTime,
	)
	if err!= nil {
        r.Logger.ERROR.Println("Error while updating reservation")
        return nil, err
    }
	res.ReservationTime = reservationTime.Format("2006-01-02")

	return &res, nil
}

func (r *ReservationRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := pb.Void{}

	query := `UPDATE reservation SET deleted_at=$1 WHERE id=$2`
	_, err := r.db.Exec(query, time.Now().Unix(), req.Id)
	if err!= nil {
        r.Logger.ERROR.Println("Error while deleting reservation")
        return nil, err
    }

	return &res, nil
}