//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/HooYa-Bigdata/rivalservice/config"
	v1 "github.com/HooYa-Bigdata/rivalservice/genproto/v1"
	"github.com/HooYa-Bigdata/rivalservice/service"
)

// InitServer Inject service's component
func InitServer(conf *config.Config) (v1.RivalServiceServer, error) {

	wire.Build(
		service.NewClient,
		service.NewServer,
	)

	return &service.Server{}, nil

}
