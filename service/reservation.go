package service

import (
	"context"
	r "reservation-service/genproto/reservation"
	st "reservation-service/storage/postgres"
)

type ReservationService struct {
	storage st.Storage
	r.UnimplementedReservationServiceServer
}

func NewReservationService(storage *st.Storage) *ReservationService {
	return &ReservationService{
		storage: *storage,
	}
}

func (s *ReservationService) Create(ctx context.Context, reservation *r.ReservationReq) (*r.Reservation, error) {
	resp, err := s.storage.ReservationS.Create(reservation)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ReservationService) Get(ctx context.Context, idReq *r.GetByIdReq) (*r.ReservationRes, error) {
	resp, err := s.storage.ReservationS.Get(idReq)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ReservationService) GetAll(ctx context.Context, allReservations *r.GetAllReservationReq) (*r.GetAllReservationRes, error) {
	reservations, err := s.storage.ReservationS.GetAll(allReservations)

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (s *ReservationService) Update(ctx context.Context, reservation *r.ReservationUpdate) (*r.Reservation, error) {
	resp, err := s.storage.ReservationS.Update(reservation)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ReservationService) Delete(ctx context.Context, idReq *r.GetByIdReq) (*r.Void, error) {
	_, err := s.storage.ReservationS.Delete(idReq)

	return nil, err
}
