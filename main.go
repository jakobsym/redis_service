package main

/*
-- starting redis using homebrew:
	-> $ brew services start redis
	-> $ redis-cli (opens redis)
	-> $ brew services stop redis
*/
import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/jakobsym/redis_service/application"
)

// storing orders made from ecommerce website via crud

// can send curl commands to server to test if working as intended
func main() {
	app := application.New()

	// returns context, if signal is created
	// context.Background() derives a new context (only use to init a context tree)
	// 'cancel()' will cancel 'ctx' and any of its child context(s) only do at end of function
	// or use defer cancel()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("Failed to start app: ", err)
	}
}
