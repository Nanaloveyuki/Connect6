package main

import (
	"log"

	"connect6/backend/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApplication("configs")
	if err != nil {
		log.Fatalf("bootstrap application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
