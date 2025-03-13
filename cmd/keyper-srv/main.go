package main

import (
	"context"
	"log"

	app "github.com/RomanAgaltsev/keyper/server/internal/app/keyper-srv"
	"github.com/RomanAgaltsev/keyper/server/internal/config"
	"github.com/RomanAgaltsev/keyper/server/internal/logger"
)

func main() {
	// Get application cofiguration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config loading: %s", err)
	}

	// Run the application
	err = app.NewApp(cfg, logger.New(cfg.Env)).Run(context.Background())
	if err != nil {
		log.Fatalf("running application : %s", err.Error())
	}
}
