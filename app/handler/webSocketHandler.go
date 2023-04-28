package handler

import (
	"c3-oc2023/models"
	"c3-oc2023/utils"
	"errors"
	"fmt"
	"io"
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
		fmt.Println(count)

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
					break
				}
				c.Logger().Error(err)
				clients.Delete(ws)
				break
			}
			ch <- msg
		}

		clients.Delete(ws)
		close(cancel)
		close(ch)
		count--
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
			switch msg.Type {
			case "pos":
				var pos models.PositionBody
				utils.MapToStruct(msg.Body.(map[string]interface{}), &pos)
				msg.Body = pos
				clients.Range(func(ws, uid interface{}) bool {
					if uid.(string) == pos.UID {
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
}
