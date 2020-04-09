package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	pb "microservices/consignment-service/proto/consignment"
	"net"
	"sync"
)

const port  = ":50051"

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu sync.Mutex
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error){
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment (ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)

	if err != nil {
		return nil, err
	}

	return &pb.Response{Created:true, Consignment:consignment}, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{repo})
	reflection.Register(s)

	log.Println("Running on port: %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serv: %v", err)
	}
}