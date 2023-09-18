// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/HooYa-Bigdata/rivalservice/config"
	"github.com/HooYa-Bigdata/rivalservice/genproto/v1"
	"github.com/HooYa-Bigdata/rivalservice/service"
)

// Injectors from wire.go:

// InitServer Inject service's component
func InitServer(conf *config.Config) (v1.RivalServiceServer, error) {
	rivalServiceClient, err := service.NewClient(conf)
	if err != nil {
		return nil, err
	}
	rivalServiceServer, err := service.NewServer(conf, rivalServiceClient)
	if err != nil {
		return nil, err
	}
	return rivalServiceServer, nil
}
