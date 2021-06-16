package app

import (
	"log"

	"github.com/joho/godotenv"
)

type App struct {
	repositories *Repositories
	Usecases     *Usecases
}

func New() (*App, error) {
	app := new(App)

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env  file \n%v", err)
	}

	repos, err := newRepositories()
	if err != nil {
		log.Fatalln(err.Error())
	}

	app.repositories = repos
	app.Usecases = newUsecases(repos)
	return app, nil
}
