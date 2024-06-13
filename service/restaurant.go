package service

import (
	"context"
	r "reservation-service/genproto/reservation"
	st "reservation-service/storage/postgres"
)

type RestaurantService struct {
	storage st.Storage
	r.UnimplementedRestaurantServiceServer
}

func NewRestaurantService(storage *st.Storage) *RestaurantService {
	return &RestaurantService{
		storage: *storage,
	}
}

func (s *RestaurantService) Create(ctx context.Context, restaurant *r.RestaurantReq) (*r.Restaurant, error) {
	resp, err := s.storage.RestaurantS.Create(restaurant)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RestaurantService) Get(ctx context.Context, idReq *r.GetByIdReq) (*r.Restaurant, error) {
	resp, err := s.storage.RestaurantS.Get(idReq)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RestaurantService) GetAll(ctx context.Context, allRestaurants *r.GetAllRestaurantReq) (*r.GetAllRestaurantRes, error) {
	restaurants, err := s.storage.RestaurantS.GetAll(allRestaurants)

	if err != nil {
		return nil, err
	}

	return restaurants, nil
}

func (s *RestaurantService) Update(ctx context.Context, restaurant *r.RestaurantUpdate) (*r.Restaurant, error) {
	resp, err := s.storage.RestaurantS.Update(restaurant)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *RestaurantService) Delete(ctx context.Context, idReq *r.GetByIdReq) (*r.Void, error) {
	_, err := s.storage.RestaurantS.Delete(idReq)

	return nil, err
}
