//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	v1 "github.com/yuri-potatoq/generic-profile/api/v1"
	"github.com/yuri-potatoq/generic-profile/enrollment"
	"github.com/yuri-potatoq/generic-profile/infra/db"
	"github.com/yuri-potatoq/generic-profile/server"
)

func injectApi(ctx context.Context) (*application, error) {
	wire.Build(
		provideApiApplication,
		db.WireSet,
		enrollment.WireSet,
		v1.WireSet,
		server.WireSet,
	)
	return &application{}, nil
}
