package application

import (
	"context"
	"fmt"
	"net/http"
)

// store application dependencies
type App struct {
	router http.Handler
}

// construrctor to init application
func New() *App {
	app := &App{
		router: loadRoutes(),
	}
	return app
}

// start application at specified route
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil // success
}
