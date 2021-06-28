package ws

import "github.com/ardafirdausr/discuss-server/internal/app"

type DiscussWebSocket struct {
	app *app.App
}

func NewDiscussWebSocket(app *app.App) *DiscussWebSocket {
	dws := &DiscussWebSocket{app: app}
	return dws
}
