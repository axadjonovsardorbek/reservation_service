package postgres

import (
	"database/sql"
	"reservation-service/config/logger"
	r "reservation-service/genproto/reservation"

	"github.com/google/uuid"
)

type MenuRepo struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewMenuRepo(db *sql.DB, logger *logger.Logger) *MenuRepo {
	return &MenuRepo{db: db, Logger: logger}
}

func(m *MenuRepo) Create(menu *r.MenuReq) (*r.Menu, error){

	id := uuid.New().String()
	res := r.Menu{}

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

func(m *MenuRepo) Get(id *r.GetByIdReq) (*r.MenuRes, error){
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

func(m *MenuRepo) GetAll(req *r.GetAllMenuReq) (*r.GetAllMenuRes, error){
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

func(m *MenuRepo) Update(menu *r.MenuUpdate) (*r.Menu, error){
	
	res := r.Menu{}

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

	row := m.db.QueryRow(query, menu.UpdateMenu.RestaurantId, menu.UpdateMenu.Name, menu.UpdateMenu.Description, menu.UpdateMenu.Price, menu.Id)

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

func(m *MenuRepo) Delete(id *r.GetByIdReq) (*r.Void, error){
	res := r.Void{}

	query := `UPDATE menu SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1 and deleted_at=0`
	_, err := m.db.Exec(query, id.Id)
	if err != nil {
		m.Logger.ERROR.Println("Error while deleting menu")
		return nil, err
	}

	return &res, nil
}
