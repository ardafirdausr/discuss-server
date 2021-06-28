package ws

import (
	"encoding/json"
	"log"

	"github.com/ardafirdausr/discuss-server/internal/app"
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/ardafirdausr/discuss-server/internal/service/pubsub"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type DiscussWebSocket struct {
	app app.App
}

func NewDiscussWebSocket(app app.App) *DiscussWebSocket {
	dws := &DiscussWebSocket{app: app}
	return dws
}

func (dws DiscussWebSocket) ListenClientMessage(wsc *WSClient) {
	for {
		_, msgContent, err := wsc.conn.ReadMessage()
		if err != nil {
			_, ok := err.(*websocket.CloseError)
			if ok {
				log.Println("connection closed by user")
				wsc.Close()
			}
			return
		}

		var message entity.CreateMessage
		if err := json.Unmarshal(msgContent, &message); err != nil {
			log.Println(err.Error())
			wsc.conn.WriteMessage(websocket.TextMessage, []byte("Failed send message"))
		}

		dws.app.Usecases.MessageUsecase.SendMessage(wsc.pubsub, message)
	}
}

func (dws DiscussWebSocket) ChatSocketHandler(c echo.Context) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal("websocket conn failed", err)
	}

	user := c.Get("user").(*entity.User)
	pbsb := pubsub.NewRedisPubSub(dws.app.Drivers.Redis)
	wsc := NewWSClient(user, conn, pbsb)

	go dws.ListenClientMessage(wsc)
}
