package main

import (
	"log"

	app "github.com/RomanAgaltsev/keyper/server/internal/app/keyper-srv"
	"github.com/RomanAgaltsev/keyper/server/internal/config"
)

func main() {
	// Get application cofiguration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config loading: %s", err)
	}

	// Run the application
	err = app.Run(cfg)
	if err != nil {
		log.Fatalf("running application : %s", err.Error())
	}
}
