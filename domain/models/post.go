package models

import "time"

type Post struct {
	Id       int64     `json:"id"`
	Parent   int64     `json:"parent"` // id родительского сообщения. 0 - корневое сообщение обсуждения
	Author   string    `json:"author"` // nickname автора сообщения
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"` // true, если сообщение было изменено
	Forum    string    `json:"forum"`    // slug форума данного сообщения
	Thread   int32     `json:"thread"`   // id ветви данного сообщения
	Created  time.Time `json:"created"`
}
