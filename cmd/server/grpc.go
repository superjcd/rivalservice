package server

import (
	"fmt"
	"log"
	"net"

	"github.com/superjcd/rivalservice/config"
	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"google.golang.org/grpc"
)

func RunGrpcServer(server v1.RivalServiceServer, cfg *config.Config) {
	grpcServer := grpc.NewServer()
	v1.RegisterRivalServiceServer(grpcServer, server)

	fmt.Println("Listening grpc server on port" + cfg.Grpc.Port)
	listen, err := net.Listen("tcp", cfg.Grpc.Port)
	if err != nil {
		panic("listen grpc tcp failed.[ERROR]=>" + err.Error())
	}

	go func() {
		if err = grpcServer.Serve(listen); err != nil {
			log.Fatal("grpc serve failed", err)
		}
	}()

	cfg.Grpc.Server = grpcServer

}
