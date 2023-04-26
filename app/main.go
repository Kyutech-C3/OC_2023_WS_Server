package main

import (
	"c3-oc2023/handler"
	"time"
)

func main() {
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	e := handler.NewHandler()

	go handler.BroadCastHandler()
	e.Logger.Fatal(e.Start(":8080"))
}
