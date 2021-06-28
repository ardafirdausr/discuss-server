package main

import (
	"log"

	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/delivery/web"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("Failed initiate the app\n%v", err)
	}

	// server
	web.Start(app)
	// ws.Start(app)
	defer app.Close()
}
