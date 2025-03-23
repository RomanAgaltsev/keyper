package main

import (
	"log"

	app "github.com/RomanAgaltsev/keyper/internal/app/keyper-cli"
	"github.com/RomanAgaltsev/keyper/internal/logger"
)

func main() {
	// Run the application
	// TODO: add cfg and log
	err := app.NewApp(logger.New("dev")).Run()
	if err != nil {
		log.Fatalf("running application : %s", err.Error())
	}
}
