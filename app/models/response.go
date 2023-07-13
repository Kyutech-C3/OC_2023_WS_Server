package models

import "time"

type Response struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type Body struct {
	UID string `json:"uid"`
}

type PositionBody struct {
	Body
	Name      string  `json:"name"`
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	Z         float32 `json:"z"`
	R         float32 `json:"r"`
	Anime     string  `json:"anime"`
	AnimeTime float32 `json:"time"`
}

type MessageBody struct {
	Body
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"at"`
}
