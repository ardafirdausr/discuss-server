package app

import (
	"context"
	"log"

	"github.com/joho/godotenv"
)

type App struct {
	Usecases     *Usecases
	Repositories *repositories
	Drivers      *drivers
}

func New() (*App, error) {
	app := new(App)

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env  file \n%v", err)
	}

	drivers, err := newDrivers()
	if err != nil {
		log.Fatalln(err)
	}

	repos := newRepositories(drivers)
	ucs := newUsecases(repos)

	app.Drivers = drivers
	app.Repositories = repos
	app.Usecases = ucs
	return app, nil
}

func (app App) Close() error {
	var rerr error

	if err := app.Drivers.Mongo.Client().Disconnect(context.TODO()); err != nil {
		log.Println("Failed to close Mongo DB connection")
		rerr = err
	}

	if err := app.Drivers.Redis.Close(); err != nil {
		log.Println("Failed to close Redis connection")
		rerr = err
	}

	return rerr
}
