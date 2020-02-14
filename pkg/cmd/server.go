package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"

	"petstore-service/config"
	"petstore-service/pkg/protocol/grpc"
	"petstore-service/pkg/protocol/rest"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type Config struct {
	Env string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.Env, "env", "", "the environment and config that used")
	flag.Parse()
	config.InitViper("./config", cfg.Env)

	if len(config.GetViper().GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", config.GetViper().GRPCPort)
	}

	if len(config.GetViper().HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", config.GetViper().HTTPPort)
	}

	// Connect to the firestore
	key := option.WithCredentialsFile(config.GetViper().Firestore.URL)
	app, err := firebase.NewApp(ctx, nil, key)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	grpc := grpc.NewGRPCServer(client)

	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, config.GetViper().GRPCPort, config.GetViper().HTTPPort)
	}()

	return grpc.RunServer(ctx, config.GetViper().GRPCPort)
}
