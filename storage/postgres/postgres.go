package postgres

import (
	"database/sql"
	"fmt"

	"reservation-service/config"
	"reservation-service/config/logger"
	"reservation-service/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db     *sql.DB
	Logger *logger.Logger
	//  RestaurantS storage.RestaurantI
	ReservationS storage.ReservationI
	// ReservationOrderS storage.ReservationOrderI
	// MenuS storage.MenuI
}

func NewPostgresStorage(config config.Config, logger *logger.Logger) (*Storage, error) {
	conn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		config.DB_HOST, config.DB_USER, config.DB_NAME, config.DB_PASSWORD, config.DB_PORT)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//  restaurant := NewRestaurantRepo(db)
	reservation := NewReservationRepo(db, logger)
	//  resOrder := NewReservationOrderRepo(db)
	//  menu := NewMenuRepo(db)

	return &Storage{
		Db: db,
		//   RestaurantS: restaurant,
		ReservationS: reservation,
		//   ReservationOrderS: resOrder,
		//   MenuS: menu,
	}, nil
}
