package handler

import (
	"c3-oc2023/models"
	"c3-oc2023/utils"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan string)

func WebSocketHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// 初期化処理
		// UIDを生成し，クライアントに送信

		// TODO: Generate Short UID
		uuid := uuid.New()
		res := &models.Response{
			Type: "init",
			Body: models.Body{
				UID: uuid.String(),
			},
		}
		bytes, err := json.Marshal(res)
		if err != nil {
			c.Logger().Error(err)
			return
		}
		if err := websocket.Message.Send(ws, string(bytes)); err != nil {
			c.Logger().Error(err)
			return
		}
		clients[ws] = uuid.String()

		go BroadCastHandler()

		// Read Message
		for {
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				delete(clients, ws)
				break
			}
			broadcast <- msg
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func BroadCastHandler() {
	for {
		msg := <-broadcast

		res := &models.Response{}
		bytes := []byte(msg)
		json.Unmarshal(bytes, &res)

		switch res.Type {
		case "pos":
			var pos models.PositionBody
			utils.MapToStruct(res.Body.(map[string]interface{}), &pos)
			for ws, uid := range clients {
				if uid == pos.UID {
					continue
				}
				if err := websocket.Message.Send(ws, msg); err != nil {
					log.Fatal(err)
					delete(clients, ws)
				}
			}

		case "mes":
			var mes models.MessageBody
			utils.MapToStruct(res.Body.(map[string]interface{}), &mes)
			for ws, uid := range clients {
				if uid == mes.UID {
					continue
				}
				if err := websocket.Message.Send(ws, msg); err != nil {
					log.Fatal(err)
					delete(clients, ws)
				}
			}
		}

	}
}
