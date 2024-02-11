package main

import (
	"context"
	"fmt"

	"github.com/jakobsym/redis_service/application"
)

// storing orders made from ecommerce website via crud

// can send curl commands to server to test if working as intended
func main() {
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Failed to start app: ", err)
	}
}
