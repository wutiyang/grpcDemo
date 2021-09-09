package main

import (
	"context"
	"log"
	"net"

	pb "github.com/wutiyang/grpcDemoconsignment-service/proto/consignment"
	"github.com/wutiyang/grpcDemoconsignment-service/repository"

	"google.golang.org/grpc"
)

const port = ":50051"

// 定义微服务
type service struct {
	repo repository.Repository
}

// service实现consignment.pb.go 中的 ShippingServiceServer 接口
// 使service作为grpc的服务端
//
// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// 接受承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	resp := &pb.Response{Created: true, Consignment: consignment}
	return resp, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	allConsignments := s.repo.GetAll()
	resp := &pb.Response{
		Created:      true,
		Consignments: allConsignments,
	}

	return resp, nil
}
func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	repo := repository.Repository{}

	// 想rpc服务器注册微服务
	// 此时会把我们自己实现的微服务service与ShippingServiceServer绑定
	pb.RegisterShippingServiceServer(server, &service{repo: repo})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
