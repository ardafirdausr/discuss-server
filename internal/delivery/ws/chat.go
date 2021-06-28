package ws

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/ardafirdausr/discuss-server/internal"
	"github.com/ardafirdausr/discuss-server/internal/entity"
	"github.com/ardafirdausr/discuss-server/internal/service/pubsub"
	"github.com/ardafirdausr/discuss-server/internal/service/token"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func (dws DiscussWebSocket) listenClientMessage(sc *socketClient) {
	_, msgContent, err := sc.conn.ReadMessage()
	if err != nil {
		_, ok := err.(*websocket.CloseError)
		if ok {
			log.Println("connection closed by user")
			sc.Close()
		}
		return
	}

	var message entity.CreateMessage
	if err := json.Unmarshal(msgContent, &message); err != nil {
		log.Println(err.Error())
		sc.conn.WriteMessage(websocket.TextMessage, []byte("Failed send message"))
	}
	message.Sender = *sc.user

	dws.app.Usecases.MessageUsecase.SendMessage(sc.pubsub, message)
}

func (dws DiscussWebSocket) listenSubscribeMessage(sc *socketClient) {
	listener := func(channel string, msg interface{}) {
		strMsg, ok := msg.(string)
		if !ok {
			err := errors.New("failed convert channel message to string")
			log.Println(err.Error())
			return
		}

		var message entity.Message
		if err := json.Unmarshal([]byte(strMsg), &message); err != nil {
			err := errors.New("failed convert channel message to message object")
			log.Println(err.Error())
			return
		}

		if message.ReceiverType == "user" && sc.user.ID != message.ReceiverID {
			err := errors.New("wrong receiver")
			log.Println(err.Error())
			return
		}

		if message.ReceiverType == "discussion" && sc.user.Discussions == nil {
			err := errors.New("wrong receiver")
			log.Println(err.Error())
			return
		}

		sc.conn.WriteMessage(websocket.TextMessage, []byte(strMsg))
	}
	sc.pubsub.Listen(listener)
}

func (dws DiscussWebSocket) authenticate(tokenizer internal.Tokenizer, strToken string) (*entity.User, error) {
	if strToken == "" {
		return nil, errors.New("token is not provided")
	}

	payload, err := tokenizer.Parse(strToken)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("invalid token")
	}

	user := &entity.User{
		ID:       payload.ID,
		Name:     payload.Name,
		Email:    payload.Email,
		ImageUrl: payload.Imageurl,
	}

	return user, nil
}

func (dws DiscussWebSocket) ChatSocketHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal("websocket conn failed", err)
	}

	strToken := c.QueryParam("token")
	JWTSecretKey := os.Getenv("JWT_SECRET_KEY")
	JWTToknizer := token.NewJWTTokenizer(JWTSecretKey)
	user, err := dws.authenticate(JWTToknizer, strToken)
	if err != nil {
		log.Println(err.Error())
		conn.WriteMessage(websocket.CloseMessage, []byte("Invalid token"))
		return conn.Close()
	}

	pbsb := pubsub.NewRedisPubSub(dws.app.Drivers.Redis)
	sc := newSocketClient(user, conn, pbsb)
	defer sc.Close()

	go dws.listenSubscribeMessage(sc)
	for {
		dws.listenClientMessage(sc)
	}
}
