// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/yuri-potatoq/generic-profile/api/v1"
	"github.com/yuri-potatoq/generic-profile/enrollment"
	"github.com/yuri-potatoq/generic-profile/infra/db"
	"github.com/yuri-potatoq/generic-profile/server"
)

// Injectors from wire.go:

func injectApi(ctx context.Context) (*application, error) {
	sqlxDB, err := db.NewDb(ctx)
	if err != nil {
		return nil, err
	}
	repository := enrollment.NewEnrollmentRepository(sqlxDB)
	service := enrollment.NewEnrollmentService(repository)
	handlerFunc := v1.GetEnrollmentHandler(service)
	serverOpts := &server.ServerOpts{
		GetEnrollmentHandler: handlerFunc,
	}
	serverServer := server.NewServer(serverOpts)
	mainApplication := provideApiApplication(sqlxDB, serverServer)
	return mainApplication, nil
}