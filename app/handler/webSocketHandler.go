package handler

import (
	"c3-oc2023/models"
	"c3-oc2023/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

var count int
var clients sync.Map
var broadcast = make(chan string, 100)

func recvHandler() {
	for {
		if len(broadcast) >= 70 {
			fmt.Println("New goroutine")
			go BroadCastHandler()
		}

		time.Sleep(1 * time.Second)
	}
}

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
		if err := websocket.JSON.Send(ws, res); err != nil {
			c.Logger().Error(err)
			return
		}
		clients.Store(ws, uuid.String())
		count++

		var wg sync.WaitGroup
		ch := make(chan models.Response, 5)
		cancel := make(chan struct{})
		wg.Add(1)
		go send(ch, cancel)

		// Read Message
		for {
			var msg models.Response
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				if errors.Is(err, io.EOF) {
					count--
					break
				}
				c.Logger().Error(err)
				clients.Delete(ws)
				count--
				break
			}
			ch <- msg
		}

		clients.Delete(ws)
		close(cancel)
		close(ch)
		wg.Wait()
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func send(ch <-chan models.Response, cancel chan struct{}) {
	for {
		select {
		case <-cancel:
			return
		default:
			msg := <-ch
			clients.Range(func(ws, uid interface{}) bool {
				if uid.(string) == msg.Body.(models.PositionBody).UID {
					return true
				}
				if err := websocket.JSON.Send(ws.(*websocket.Conn), msg); err != nil {
					if errors.Is(err, syscall.EPIPE) {
						clients.Delete(ws.(*websocket.Conn))
						return true
					} else {
						return false
					}
				}
				return true
			})
		}
	}
}

func BroadCastHandler() {
	for {
		msg := <-broadcast

		res := &models.Response{}
		bytes := []byte(msg)
		json.Unmarshal(bytes, &res)
		// fmt.Println(res)

		switch res.Type {
		case "pos":
			var pos models.PositionBody
			utils.MapToStruct(res.Body.(map[string]interface{}), &pos)
			// fmt.Println(pos)
			clients.Range(func(ws, uid any) bool {
				if uid == pos.UID {
					return true
				}
				if err := websocket.Message.Send(ws.(*websocket.Conn), msg); err != nil {
					log.Fatal(err)
					ws.(*websocket.Conn).Close()
					clients.Delete(ws)
					count--
				}
				return true
			})

			// case "mes":
			// 	var mes models.MessageBody
			// 	utils.MapToStruct(res.Body.(map[string]interface{}), &mes)
			// 	for ws, uid := range clients {
			// 		if uid == mes.UID {
			// 			continue
			// 		}
			// 		if err := websocket.Message.Send(ws, msg); err != nil {
			// 			log.Fatal(err)
			// 			delete(clients, ws)
			// 		}
			// 	}
		}

	}
}
