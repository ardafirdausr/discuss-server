package ws

import (
	"net/http"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.URL.Hostname() == "localhost"
	},
}

type WSClient struct {
	user   *entity.User
	conn   *websocket.Conn
	pubsub internal.PubSub
}

func NewWSClient(user *entity.User, conn *websocket.Conn, pubsub internal.PubSub) *WSClient {
	return &WSClient{user: user, conn: conn, pubsub: pubsub}
}

func (wsc WSClient) Close() {
	wsc.conn.Close()
	wsc.pubsub.Close()
}
