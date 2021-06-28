package ws

import (
	"log"
	"net/http"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type socketClient struct {
	user   *entity.User
	conn   *websocket.Conn
	pubsub internal.PubSub
}

func newSocketClient(user *entity.User, conn *websocket.Conn, pubsub internal.PubSub) *socketClient {
	return &socketClient{user: user, conn: conn, pubsub: pubsub}
}

func (wsc socketClient) Close() error {
	var rerr error

	if err := wsc.conn.Close(); err != nil {
		log.Println("Failed to close websocket connection")
		rerr = err
	}

	if err := wsc.pubsub.Close(); err != nil {
		log.Println("Failed to close subscription connection")
		rerr = err
	}

	return rerr
}
