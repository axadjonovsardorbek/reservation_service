package service

import (
	r "reservation-service/genproto/reservation"
	st "reservation-service/storage/postgres"
)

type MenuService struct {
	storage st.Storage
	r.UnimplementedMenuServiceServer
}

func NewMenuService(storage *st.Storage) *MenuService {
	return &MenuService{
		storage: *storage,
	}
}

// func (s *MenuService) Create(ctx context.Context, menu *r.MenuReq) (*r.Menu, error) {
// 	resp, err := s.storage.MenuS.Create(menu)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (s *MenuService) Get(ctx context.Context, idReq *r.GetByIdReq) (*r.MenuRes, error) {
// 	resp, err := s.storage.MenuS.Get(idReq)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (s *MenuService) GetAll(ctx context.Context, allMenus *r.GetAllMenuReq) (*r.GetAllMenuRes, error) {
// 	items, err := s.storage.MenuS.GetAll(allMenus)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return items, nil
// }

// func (s *MenuService) Update(ctx context.Context, menu *r.MenuUpdate) (*r.Menu, error) {
// 	resp, err := s.storage.MenuS.Update(menu)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (s *MenuService) Delete(ctx context.Context, idReq *r.GetByIdReq) (*r.Void, error) {
// 	_, err := s.storage.MenuS.Delete(idReq)

// 	return nil, err
// }
