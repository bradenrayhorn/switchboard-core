package grpc

import (
	"fmt"
	"github.com/bradenrayhorn/switchboard-protos/groups"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s *Server) Start() {
	requestedPort := viper.GetString("grpc_port")
	log.Printf("starting grpc server on port %s", requestedPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", requestedPort))
	if err != nil {
		log.Fatalf("failed to bind grpc port %s: %v", requestedPort, err)
	}

	grpcServer := grpc.NewServer()

	groups.RegisterGroupServiceServer(grpcServer, &GroupsServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server: %s", err)
	}
}
