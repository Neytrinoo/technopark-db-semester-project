package models

import "time"

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`    // nickname полоьзователя-владельца форума
	Slug    string `json:"slug"`    // уникальное поле, человекопонятный урл форума
	Posts   int64  `json:"posts"`   // кол-во постов в форуме
	Threads int32  `json:"threads"` // кол-во ветвей обсуждения в форуме
}

type ForumCreate struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type GetForumUsers struct {
	Slug  string `json:"slug"`
	Limit int32  `json:"limit"` // default 100
	Since string `json:"since"` // nickname пользователя, с которого будем выводить результат
	Desc  bool   `json:"desc"`
}

type GetForumThreads struct {
	Limit int32     `json:"limit,omitempty"` // default 100
	Since time.Time `json:"since,omitempty"` // nickname пользователя, с которого будем выводить результат
	Desc  bool      `json:"desc,omitempty"`
}
