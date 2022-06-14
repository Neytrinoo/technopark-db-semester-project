package models

import "time"

type Forum struct {
	Title   string `json:"title" db:"title"`
	User    string `json:"user" db:"user"`       // nickname полоьзователя-владельца форума
	Slug    string `json:"slug" db:"slug"`       // уникальное поле, человекопонятный урл форума
	Posts   int64  `json:"posts" db:"posts"`     // кол-во постов в форуме
	Threads int32  `json:"threads" db:"threads"` // кол-во ветвей обсуждения в форуме
}

type ForumCreate struct {
	Title string `json:"title" db:"title"`
	User  string `json:"user" db:"user"`
	Slug  string `json:"slug" db:"slug"`
}

type GetForumUsers struct {
	Slug  string `json:"slug" db:"slug"`
	Limit int32  `json:"limit" db:"limit"` // default 100
	Since string `json:"since" db:"since"` // nickname пользователя, с которого будем выводить результат
	Desc  bool   `json:"desc" db:"desc"`
}

type GetForumThreads struct {
	Limit int32     `json:"limit,omitempty"` // default 100
	Since time.Time `json:"since,omitempty"` // nickname пользователя, с которого будем выводить результат
	Desc  bool      `json:"desc,omitempty"`
}
