package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/yuri-potatoq/generic-profile/server"
)

type application struct {
	api struct {
		server *server.Server
	}

	pool *sqlx.DB
}

func (app *application) initApi() error {
	return app.api.server.Start()
}

func (app *application) closeDb() error {
	return app.pool.Close()
}

func provideApiApplication(
	pool *sqlx.DB,
	server *server.Server,
) *application {
	app := new(application)
	app.pool = pool
	app.api.server = server
	return app
}
