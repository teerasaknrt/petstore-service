package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	petstoreAPI "petstore-service/pkg/api/v1"
	petstoreService "petstore-service/pkg/service"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	petstore petstoreAPI.PetStoreServiceServer
}

func NewGRPCServer(db *firestore.Client) *GRPCServer {
	return &GRPCServer{
		petstore: petstoreService.NewPetStoreServiceServer(db),
	}
}

// RunServer runs gRPC service to publish Auth service
func (s *GRPCServer) RunServer(ctx context.Context, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	// gRPC server statup otions
	opts := []grpc.ServerOption{}

	// register service
	server := grpc.NewServer(opts...)
	petstoreAPI.RegisterPetStoreServiceServer(server, s.petstore)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")
			//postgres.DisconnectPostgres(ctx, s.Postgres)

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server..." + port)
	return server.Serve(listen)
}
