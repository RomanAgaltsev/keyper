package main

import (
	"log"

	app "github.com/RomanAgaltsev/keyper/internal/app/keyper-srv"
	"github.com/RomanAgaltsev/keyper/internal/config"
	"github.com/RomanAgaltsev/keyper/internal/logger"
)

func main() {
	// Get application cofiguration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config loading: %s", err)
	}

	// Run the application
	err = app.NewApp(cfg, logger.New(cfg.Env)).Run()
	if err != nil {
		log.Fatalf("running application : %s", err.Error())
	}
}
