package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/freekobie/kora/handler"
	"github.com/freekobie/kora/service"
)

type application struct {
	handler *handler.Handler
	server  *http.Server
	fileService *service.FileService
}

func newApplication(handler *handler.Handler, address string, fileService *service.FileService) *application {
	server := http.Server{
		Addr: fmt.Sprintf(":%s", address),
	}

	return &application{
		handler: handler,
		server:  &server,
		fileService: fileService,
	}
}

func (app *application) start() error {

	app.server.Handler = app.routes()

	return app.server.ListenAndServe()
}

func (app *application) shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
