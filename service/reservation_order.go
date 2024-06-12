package service

import (
	r "reservation-service/genproto/reservation"
	st "reservation-service/storage/postgres"
)

type ReservationOrderService struct {
	storage st.Storage
	r.UnimplementedReservationOrderServiceServer
}

func NewReservationOrderService(storage *st.Storage) *ReservationOrderService {
	return &ReservationOrderService{
		storage: *storage,
	}
}

// func (s *ReservationOrderService) Create(ctx context.Context, order *r.ReservationOrderReq) (*r.ReservationOrderRes, error) {
// 	resp, err := s.storage.ReservationOrderS.Create(order)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (s *ReservationOrderService) Get(ctx context.Context, idReq *r.GetByIdReq) (*r.ReservationOrderRes, error) {
// 	resp, err := s.storage.ReservationOrderS.Get(idReq)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (s *ReservationOrderService) GetAll(ctx context.Context, allOrders *r.GetAllReservationOrderReq)(*r.GetAllReservationOrderRes, error){
// 	orders, err := s.storage.ReservationOrderS.GetAll(allOrders)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

// func (s *ReservationOrderService) Update(ctx context.Context, reservation *r.ReservationOrderUpdate)(*r.ReservationOrderRes, error){
// 	resp, err := s.storage.ReservationOrderS.Update(reservation)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (s *ReservationOrderService) Delete(ctx context.Context, idReq *r.GetByIdReq)(*r.Void, error){
// 	_, err := s.storage.ReservationOrderS.Delete(idReq)

// 	return nil, err
// }
