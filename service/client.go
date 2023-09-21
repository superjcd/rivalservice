package service

import (
	"context"
	"time"

	"github.com/superjcd/rivalservice/config"
	v1 "github.com/superjcd/rivalservice/genproto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient New service's client
func NewClient(conf *config.Config) (v1.RivalServiceClient, error) {

	serverAddress := conf.Grpc.Host + conf.Grpc.Port
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := v1.NewRivalServiceClient(conn)
	return client, nil

}
