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
	Name string `json:"name"`
	X    string `json:"x"`
	Y    string `json:"y"`
}

type MessageBody struct {
	Body
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"at"`
}
