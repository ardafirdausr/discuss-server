package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/ardafirdausr/discuss-server/internal/service/pubsub"
	"github.com/ardafirdausr/discuss-server/internal/service/token"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func (dws DiscussWebSocket) listenClientMessage(sc *socketClient, quit chan<- bool) {
	for {
		_, msgContent, err := sc.conn.ReadMessage()
		if _, ok := err.(*websocket.CloseError); ok {
			log.Println(err.Error())
			quit <- true
			return
		}

		if err != nil {
			log.Println(err.Error())
			sc.conn.WriteMessage(websocket.TextMessage, []byte("Failed send message"))
		}

		var message entity.CreateMessage
		if err := json.Unmarshal(msgContent, &message); err != nil {
			log.Println(err.Error())
			sc.conn.WriteMessage(websocket.TextMessage, []byte("Failed send message"))
		}
		message.Sender = *sc.user

		if _, err := dws.app.Usecases.MessageUsecase.SendMessage(sc.pubsub, message); err != nil {
			log.Println(err.Error())
			sc.conn.WriteMessage(websocket.TextMessage, []byte("Failed send message"))
		}
	}
}

func (dws DiscussWebSocket) listenSubscribeMessage(sc *socketClient) {
	listener := func(channel string, msg string) {
		var message entity.Message
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			err := errors.New("failed convert channel message to message object")
			log.Println(err.Error())
			return
		}

		if message.Sender.ID == sc.user.ID {
			return
		}

		sc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
	sc.pubsub.Listen(listener)
}

func (dws DiscussWebSocket) ChatSocketHandler(c echo.Context) error {
	wsLogger := log.New(log.Writer(), "SOCKET /ws/chat ", log.Ldate|log.Ltime|log.Lmsgprefix)
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal("websocket conn failed", err)
	}

	strToken := c.QueryParam("token")
	JWTSecretKey := os.Getenv("JWT_SECRET_KEY")
	JWTToknizer := token.NewJWTTokenizer(JWTSecretKey)
	user, err := dws.app.Usecases.AuthUsecase.GetUserFromToken(strToken, JWTToknizer)
	if err != nil {
		log.Println(err.Error())
		conn.WriteMessage(websocket.CloseMessage, []byte("Invalid token"))
		return conn.Close()
	}

	pbsb := pubsub.NewRedisPubSub(dws.app.Drivers.Redis)
	sc := newSocketClient(user, conn, pbsb)
	defer func() {
		sc.Close()
		wsLogger.Printf("%s disconnected from the chat socket \n", user.Email)
	}()

	var discussionChannels []string
	for _, discussion := range user.Discussions {
		discussionChannel := fmt.Sprintf("%s/%v", entity.MessageReceiverDiscussion, discussion.ID)
		discussionChannels = append(discussionChannels, discussionChannel)
	}
	sc.pubsub.Subscribe(discussionChannels...)

	quit := make(chan bool)
	defer close(quit)
	wsLogger.Printf("%s connected to the chat socket \n", user.Email)
	go dws.listenSubscribeMessage(sc)
	go dws.listenClientMessage(sc, quit)

	<-quit
	return nil
}
