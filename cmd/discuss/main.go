package main

import (
	"fmt"
	"log"

	"github.com/ardafirdausr/discuss-server/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("Failed initiate the app\n%v", err)
	}

	fmt.Println(app)
}
