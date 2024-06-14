package postgres

import (
	"database/sql"
	"fmt"
	"reservation-service/genproto/reservation"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func NewTestReservation(t *testing.T) *ReservationRepo {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		"postgres",
		"feruza1727",
		"localhost",
		"5432",
		"reservation")

	db, err := sql.Open("postgres", connString)
	if err != nil {
		t.Fatal("failed to open connection", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal("failed to open connection", err)
	}
	return &ReservationRepo{db: db}
}

func TestCreateReservation(t *testing.T) {
	stgReservation := NewTestReservation(t)
	defer stgReservation.db.Close()
	testReservation := reservation.ReservationReq{
		UserId:          "a03150a8-598b-4a5e-8283-ccadcdc4a005",
		RestaurantId:    "d13cc5a3-0be1-4e5c-87be-53e509d4edf9",
		ReservationTime: "2023-12-12 12:12:12",
		Status:          "bron qilingan",
	}

	reservationRes, err := stgReservation.Create(&testReservation)

	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}

	assert.NotNil(t, reservationRes)
	assert.Equal(t, testReservation.UserId, reservationRes.UserId)
	assert.Equal(t, testReservation.RestaurantId, reservationRes.RestaurantId)
	assert.NotEqual(t, testReservation.ReservationTime, reservationRes.ReservationTime)
	assert.Equal(t, "bron qilingan", reservationRes.Status)
}

func TestGetReservation(t *testing.T) {
	// Test bazasini yaratish
	stgReservation := NewTestReservation(t)
	defer stgReservation.db.Close()

	// Yangi test rezervatsiyani yaratish
	testReservation := reservation.ReservationReq{
		UserId:          "a03150a8-598b-4a5e-8283-ccadcdc4a005",
		RestaurantId:    "ba6ca424-395c-4b9e-a0a1-d01197084e0e",
		ReservationTime: "2023-12-12 12:12:12",
		Status:          "bron qilingan",
	}

	createdReservation, err := stgReservation.Create(&testReservation)
	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}

	// Get metodini tekshirish
	reservationGet, err := stgReservation.Get(&reservation.GetByIdReq{Id: createdReservation.Id})
	if err != nil {
		t.Fatalf("failed to get reservation: %v", err)
	}

	assert.NotNil(t, reservationGet)
	assert.Equal(t, createdReservation.Id, reservationGet.Id)
	assert.NotEqual(t, createdReservation.UserId, reservationGet.User.Id)
	assert.Equal(t, createdReservation.RestaurantId, reservationGet.Restaurant.Id)
	assert.NotEqual(t, createdReservation.ReservationTime, reservationGet.ReservationTime)
	assert.Equal(t, createdReservation.Status, reservationGet.Status)
}

func TestGetAllReservation(t *testing.T) {
	stgReservation := NewTestReservation(t)
	defer stgReservation.db.Close()

	// Yangi test rezervatsiyani yaratish
	for i := 0; i < 10; i++ {
		testReservation := reservation.ReservationReq{
			UserId:          "a03150a8-598b-4a5e-8283-ccadcdc4a005",
			RestaurantId:    "ba6ca424-395c-4b9e-a0a1-d01197084e0e",
			ReservationTime: "2023-12-12 12:12:12",
			Status:          "bron qilingan",
		}
		_, err := stgReservation.Create(&testReservation)
		if err != nil {
			t.Fatalf("failed to create reservation: %v", err)
		}
	}

	// Get metodini tekshirish
	reservationGet, err := stgReservation.GetAll(&reservation.GetAllReservationReq{
		UserId: "a03150a8-598b-4a5e-8283-ccadcdc4a005",
		Filter: &reservation.Filter{
			Limit:  10,
			Offset: 1,
		},
	})
	if err != nil {
		t.Fatalf("failed to get reservation: %v", err)
	}
	assert.GreaterOrEqual(t, 111, len(reservationGet.Reservation))

}

func TestUpdateReservation(t *testing.T) {
	stgReservation := NewTestReservation(t)
	defer stgReservation.db.Close()
	testReservation := reservation.ReservationReq{
		UserId:          "a03150a8-598b-4a5e-8283-ccadcdc4a005",
		RestaurantId:    "ba6ca424-395c-4b9e-a0a1-d01197084e0e",
		ReservationTime: "2023-12-12 12:12:12",
		Status:          "bron qilingan",
	}
	createdReservation, err := stgReservation.Create(&testReservation)
	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}
	updatedReservation := reservation.ReservationReq{
		UserId:          "a03150a8-598b-4a5e-8283-ccadcdc4a005",
		RestaurantId:    "ba6ca424-395c-4b9e-a0a1-d01197084e0e",
		ReservationTime: "2023-12-12 12:12:12",
		Status:          "bron qilingan",
	}
	_, err = stgReservation.Update(&reservation.ReservationUpdate{
		Id: &reservation.GetByIdReq{
			Id: createdReservation.Id,
		},
		UpdateReservation: &updatedReservation,
	})
	if err != nil {
		t.Fatalf("failed to update reservation: %v", err)
	}

	// Get metodini tekshirish
	reservationGet, err := stgReservation.Get(&reservation.GetByIdReq{
		Id: createdReservation.Id,
	})
	if err != nil {
		t.Fatalf("failed to get reservation: %v", err)
	}
	assert.Equal(t, createdReservation.Id, reservationGet.Id)
	assert.NotEqual(t, createdReservation.UserId, reservationGet.User.Id)
	assert.Equal(t, createdReservation.RestaurantId, reservationGet.Restaurant.Id)
	assert.Equal(t, updatedReservation.ReservationTime, reservationGet.ReservationTime)
	assert.Equal(t, updatedReservation.Status, reservationGet.Status)
}

func TestDeleteReservation(t *testing.T) {
	stgReservation := NewTestReservation(t)
	defer stgReservation.db.Close()
	testReservation := reservation.ReservationReq{
		UserId:          "a03150a8-598b-4a5e-8283-ccadcdc4a005",
		RestaurantId:    "ba6ca424-395c-4b9e-a0a1-d01197084e0e",
		ReservationTime: "2023-12-12 12:12:12",
		Status:          "bron qilingan",
	}
	createdReservation, err := stgReservation.Create(&testReservation)
	if err != nil {
		t.Fatalf("failed to create reservation: %v", err)
	}
	_, err = stgReservation.Delete(&reservation.GetByIdReq{
		Id: createdReservation.Id,
	})
	if err != nil {
		t.Fatalf("failed to delete reservation: %v", err)
	}
}
