package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

// store application dependencies

type App struct {
	router http.Handler
	rdb    *redis.Client
}

// construrctor to init application
func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}
	return app
}

// start application at specified route
func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	// check if db online
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	// close db using 'defer'
	defer func() {
		if err := a.rdb.Close(); err != nil {
			err := fmt.Errorf("failed to close redis: %w", err)
			fmt.Println(err)
		}
	}()

	fmt.Println("Starting server.")

	ch := make(chan error, 1)
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()
	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		// give time for anything to finish before shutdown
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}

	// can check if channel was closed below
	// err, open := <-ch
	/*
		if !open {
			// channel is closed
		}
	*/
}
