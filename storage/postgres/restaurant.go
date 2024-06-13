package postgres

import (
	"database/sql"
	"fmt"
	"reservation-service/config/logger"
	pb "reservation-service/genproto/reservation"

	"github.com/google/uuid"
)

type RestaurantRepo struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewRestaurantRepo(db *sql.DB, logger *logger.Logger) *RestaurantRepo {
	return &RestaurantRepo{db: db, Logger: logger}
}

func (r *RestaurantRepo) Create(restaurant *pb.RestaurantReq) (*pb.Restaurant, error) {

	id := uuid.New().String()
	res := pb.Restaurant{}

	query := `
	INSERT INTO restaurants (
		id,
		name,
		address,
		phone_number,
		description
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING 
		id,
		name,
		address,
		phone_number,
		description
	`

	row := r.db.QueryRow(query, id, restaurant.Name, restaurant.Address, restaurant.PhoneNumber, restaurant.Description)

	err := row.Scan(
		&res.Id,
		&res.Name,
		&res.Address,
		&res.PhoneNumber,
		&res.Description,
	)

	if err != nil {
		r.Logger.ERROR.Println("Error while creating restaurant")
		return nil, err
	}

	r.Logger.INFO.Println("Successfully created restaurant")

	return &res, nil
}

func (r *RestaurantRepo) Get(id *pb.GetByIdReq) (*pb.Restaurant, error) {
	var res pb.Restaurant

	query := `SELECT id, name, address, phone_number, description FROM restaurants WHERE id = $1 AND deleted_at=0`

	row := r.db.QueryRow(query, id.Id)
	err := row.Scan(
		&res.Id,
		&res.Name,
		&res.Address,
		&res.PhoneNumber,
		&res.Description,
	)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting restaurant by id : ", err)
		return nil, err
	}

	return &res, nil
}

func (r *RestaurantRepo) GetAll(req *pb.GetAllRestaurantReq) (*pb.GetAllRestaurantRes, error) {

	res := &pb.GetAllRestaurantRes{
		Restaurant: []*pb.Restaurant{},
	}

	query := `SELECT id, name, address, phone_number, description FROM restaurants WHERE deleted_at=0 LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, req.Filter.Limit, req.Filter.Offset)
	if err != nil {
		r.Logger.ERROR.Println("Error while getting restaurants : ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		restaurant := pb.Restaurant{}
		err := rows.Scan(
			&restaurant.Id,
			&restaurant.Name,
			&restaurant.Address,
			&restaurant.PhoneNumber,
			&restaurant.Description,
		)
		if err != nil {
			r.Logger.ERROR.Println("Error while getting restaurant : ", err)
			return nil, err
		}
		res.Restaurant = append(res.Restaurant, &restaurant)
	}
	err = rows.Err()
	if err != nil {
		r.Logger.ERROR.Println("Error while getting restaurants : ", err)
		return nil, err
	}

	return res, nil
}

func (r *RestaurantRepo) Update(restaurant *pb.RestaurantUpdate) (*pb.Restaurant, error) {

	if r.db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }
    if r.Logger == nil {
        return nil, fmt.Errorf("logger is nil")
    }


	res := pb.Restaurant{}

	query := `
	UPDATE restaurants SET
		name=$1,
		address=$2,
		phone_number=$3,
		description=$4,
		updated_at=now()
	WHERE 
		id=$5
	AND 
		deleted_at = 0
	RETURNING
		id,
		name,
		address,
		phone_number,
		description
	`

	row := r.db.QueryRow(query, restaurant.UpdateRestaurant.Name, restaurant.UpdateRestaurant.Address, restaurant.UpdateRestaurant.PhoneNumber, restaurant.UpdateRestaurant.Description, restaurant.Id.Id)

	err := row.Scan(
		&res.Id,
		&res.Name,
		&res.Address,
		&res.PhoneNumber,
		&res.Description,
	)

	if err != nil {
		r.Logger.ERROR.Println("Error while updating restaurant")
		return nil, err
	}

	r.Logger.INFO.Println("Successfully updated restaurant")

	return &res, nil
}

func (r *RestaurantRepo) Delete(id *pb.GetByIdReq) (*pb.Void, error) {

	res := pb.Void{}

	query := `UPDATE restaurants SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1 and deleted_at=0`
	_, err := r.db.Exec(query, id.Id)
	if err != nil {
		r.Logger.ERROR.Println("Error while deleting restaurant")
		return nil, err
	}

	r.Logger.INFO.Println("Successfully deleted restaurant")
	return &res, nil
}
