package app

type App struct {
	repositories *Repositories
	usecases     *UseCases
}

func New() *App {
	return &App{
		repositories: newRepositories(),
		usecases:     newUseCases(),
	}
}
