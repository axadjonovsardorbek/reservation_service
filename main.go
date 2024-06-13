package main

import (
	"log"
	"net"
	"path/filepath"
	"runtime"

	cf "reservation-service/config"
	"reservation-service/config/logger"

	pb "reservation-service/genproto/reservation"
	service "reservation-service/service"
	"reservation-service/storage/postgres"

	"google.golang.org/grpc"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	config := cf.Load()
	logger := logger.NewLogger(basepath, config.LOG_PATH) // Don't forget to change your log path
	em := cf.NewErrorManager(logger)
	db, err := postgres.NewPostgresStorage(config, logger)
	em.CheckErr(err)
	defer db.Db.Close()

	listener, err := net.Listen("tcp", config.RESERVATION_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterReservationServiceServer(s, service.NewReservationService(db))
	pb.RegisterRestaurantServiceServer(s, service.NewRestaurantService(db))
	pb.RegisterMenuServiceServer(s, service.NewMenuService(db))
	pb.RegisterReservationOrderServiceServer(s, service.NewReservationOrderService(db))

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
