package postgres

import (
	"database/sql"
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

func(r *RestaurantRepo) Create(restaurant *pb.RestaurantReq) (*pb.Restaurant, error){

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

func(r *RestaurantRepo) Get(id *pb.GetByIdReq) (*pb.Restaurant, error){
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

func(r *RestaurantRepo) GetAll(req *pb.GetAllRestaurantReq) (*pb.GetAllRestaurantRes, error){
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

func(r *RestaurantRepo) Update(restaurant *pb.RestaurantUpdate) (*pb.Restaurant, error){

	res := pb.Restaurant{}

	query := `
	UPDATE restaurants SET
		name=$1,
		address=$2,
		phone_number=$3,
		description=$4
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

	row := r.db.QueryRow(query, restaurant.UpdateRestaurant.Name, restaurant.UpdateRestaurant.Address, restaurant.UpdateRestaurant.PhoneNumber, restaurant.UpdateRestaurant.Description, restaurant.Id)

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

func(r *RestaurantRepo) Delete(id *pb.GetByIdReq) (*pb.Void, error){
	
	res := pb.Void{}

	query := `UPDATE restaurants SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1 and deleted_at=0`
	_, err := r.db.Exec(query, id.Id)
	if err != nil {
		r.Logger.ERROR.Println("Error while deleting restaurant")
		return nil, err
	}

	return &res, nil
}