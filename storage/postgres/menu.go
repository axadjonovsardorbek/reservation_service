package postgres

import (
	"database/sql"
	"fmt"
	"reservation-service/config/logger"
	pb "reservation-service/genproto/reservation"

	"github.com/google/uuid"
)

type MenuRepo struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewMenuRepo(db *sql.DB, logger *logger.Logger) *MenuRepo {
	return &MenuRepo{db: db, Logger: logger}
}

func (m *MenuRepo) Create(menu *pb.MenuReq) (*pb.Menu, error) {

	id := uuid.New().String()
	res := pb.Menu{}

	query := `
	INSERT INTO menu (
		id,
		restaurant_id,
		name,
		description,
		price
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING 
		id,
		restaurant_id,
		name,
		description,
		price
	`

	row := m.db.QueryRow(query, id, menu.RestaurantId, menu.Name, menu.Description, menu.Price)

	err := row.Scan(
		&res.Id,
		&res.RestaurantId,
		&res.Name,
		&res.Description,
		&res.Price,
	)

	if err != nil {
		m.Logger.ERROR.Println("Error while creating menu")
		return nil, err
	}

	m.Logger.INFO.Println("Successfully created menu")

	return &res, nil
}

func (m *MenuRepo) Get(id *pb.GetByIdReq) (*pb.MenuRes, error) {
	res := pb.MenuRes{
		Restaurant: &pb.Restaurant{},
	}

	query := `SELECT 
                m.id, 
                r.id as restaurant_id,
                r.name as restaurant_name,
                r.description as restaurant_description, 
                r.address,
                r.phone_number,
                m.name, 
                m.description, 
                m.price 
            FROM menu m 
            JOIN restaurants r ON m.restaurant_id = r.id
            WHERE m.id = $1 AND m.deleted_at=0`

	row := m.db.QueryRow(query, id.Id)

	err := row.Scan(
		&res.Id,
		&res.Restaurant.Id,
		&res.Restaurant.Name,
		&res.Restaurant.Description,
		&res.Restaurant.Address,
		&res.Restaurant.PhoneNumber,
		&res.Name,
		&res.Description,
		&res.Price,
	)
	if err != nil {
		m.Logger.ERROR.Println("Error while getting menu by id : ", err)
		return nil, err
	}

	return &res, nil
}

func (m *MenuRepo) GetAll(req *pb.GetAllMenuReq) (*pb.GetAllMenuRes, error) {

	res := &pb.GetAllMenuRes{
		Menu: []*pb.MenuRes{},
	}

	query := `SELECT 
				m.id, 
				r.id as restaurant_id,
                r.name as restaurant_name,
                r.description as restaurant_description, 
                r.address,
                r.phone_number,
				m.name, 
                m.description, 
                m.price
			FROM menu m
			JOIN restaurants r ON m.restaurant_id = r.id
			WHERE m.deleted_at=0 LIMIT $1 OFFSET $2`

	rows, err := m.db.Query(query, req.Filter.Limit, req.Filter.Offset)
	if err != nil {
		m.Logger.ERROR.Println("Error while getting menus : ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		menu := pb.MenuRes{
			Restaurant: &pb.Restaurant{},
		}
		err := rows.Scan(
			&menu.Id,
			&menu.Restaurant.Id,
			&menu.Restaurant.Name,
			&menu.Restaurant.Description,
			&menu.Restaurant.Address,
			&menu.Restaurant.PhoneNumber,
			&menu.Name,
			&menu.Description,
			&menu.Price,
		)
		if err != nil {
			m.Logger.ERROR.Println("Error while getting menus : ", err)
			return nil, err
		}
		res.Menu = append(res.Menu, &menu)
	}

	return res, nil
}

func (m *MenuRepo) Update(menu *pb.MenuUpdate) (*pb.Menu, error) {

	res := pb.Menu{}

	query := `
	UPDATE menu SET
		restaurant_id=$1,
		name=$2,
		description=$3,
		price=$4
	WHERE 
		id=$5
	AND 
		deleted_at = 0
	RETURNING
		id,
		restaurant_id,
		name,
		description,
		price
	`

	row := m.db.QueryRow(query, menu.UpdateMenu.RestaurantId, menu.UpdateMenu.Name, menu.UpdateMenu.Description, menu.UpdateMenu.Price, menu.Id.Id)

	err := row.Scan(
		&res.Id,
		&res.RestaurantId,
		&res.Name,
		&res.Description,
		&res.Price,
	)

	if err != nil {
		m.Logger.ERROR.Println("Error while updating menu")
		return nil, err
	}

	m.Logger.INFO.Println("Successfully updated menu")

	return &res, nil
}

func (m *MenuRepo) Delete(id *pb.GetByIdReq) (*pb.Void, error) {
	res := pb.Void{}

	query := `UPDATE menu SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1 and deleted_at=0`
	ress, err := m.db.Exec(query, id.Id)
	if err != nil {
		m.Logger.ERROR.Println("Error while deleting menu")
		return nil, err
	}

	if r, err := ress.RowsAffected(); r == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("restaurant with id %s not found", id.Id)
	}

	return &res, nil
}
