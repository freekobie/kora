package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/freekobie/kora/handler"
)

type application struct {
	handler *handler.Handler
	server  *http.Server
}

func newApplication(handler *handler.Handler, address string) *application {
	server := http.Server{
		Addr: fmt.Sprintf(":%s", address),
	}

	return &application{
		handler: handler,
		server:  &server,
	}
}

func (app *application) start() error {

	app.server.Handler = app.routes()

	return app.server.ListenAndServe()
}

func (app *application) shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
