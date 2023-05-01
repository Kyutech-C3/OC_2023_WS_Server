package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/websocket"
)

const (
	url    = "wss://oc2023-ws.compositecomputer.club/ws"
	origin = "http://localhost"
	fps    = 30
)

var counter int

type Cancel struct{}

type (
	Response struct {
		Type string `json:"type"`
		Body Body   `json:"body"`
	}
	Body struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
		X    string `json:"x"`
		Y    string `json:"y"`
	}
)

func main() {
	chs := []chan Cancel{}
	var num int
	fmt.Scan(&num)

	end := make(chan Cancel)
	go Interval(1*time.Second, end, func() {
		fmt.Println(counter)
	})

	for num != 0 {
		if num > 0 {
			// Open
			for num > 0 {
				ch := make(chan Cancel)
				chs = append(chs, ch)
				go newClient(ch)
				num--
			}
		} else {
			// Close
			num *= -1
			if len(chs) < num {
				num = len(chs)
			}
			for _, ch := range chs[:num] {
				close(ch)
			}

			chs = chs[num:]
		}

		fmt.Printf("Current num of clients : %d\n", len(chs))
		fmt.Scan(&num)
	}

	for _, ch := range chs {
		close(ch)
	}
	close(end)
}

func newClient(cancel chan Cancel) {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	// var init struct {
	// 	Type string `json:"type"`
	// 	Body struct {
	// 		UID string `json:"uid"`
	// 	} `json:"body"`
	// }
	// if err := websocket.JSON.Receive(ws, &init); err != nil {
	// 	panic(err)
	// }

	go recvMsg(ws)

	for {
		select {
		case <-cancel:
			// 中断通知が来た
			return
		default:
			msg := &Response{
				Type: "pos",
				Body: Body{
					UID: "test",
					// UID:  init.Body.UID,
					Name: "test",
					X:    "0",
					Y:    "0",
				},
			}
			websocket.JSON.Send(ws, msg)
			time.Sleep(1 * time.Second / fps)
		}
	}
}

func recvMsg(ws *websocket.Conn) {
	var rcvMsg Response
	var count int
	ch := make(chan Cancel)
	go Interval(1*time.Second, ch, func() {
		counter = count
		count = 0
	})

	for {
		// var msg string
		// if err := websocket.Message.Receive(ws, &msg); err != nil {
		if err := websocket.JSON.Receive(ws, &rcvMsg); err != nil {
			if errors.Is(err, io.EOF) {
				close(ch)
				ws.Close()
			}
		}
		// fmt.Println(msg)
		count++
	}
}

func Interval(sec time.Duration, ch chan Cancel, f func()) {
	for {
		select {
		case <-ch:
			return
		default:
			f()
			time.Sleep(sec)
		}
	}
}
