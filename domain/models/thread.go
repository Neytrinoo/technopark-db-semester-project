package models

import "time"

type Thread struct {
	Id      int32     `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"` // nickname пользователя, создавшего ветку
	Forum   string    `json:"forum"`  // slug форума, в котором создана данная ветка
	Message string    `json:"message"`
	Votes   int32     `json:"votes"`          // кол-во голосов за данное сообщение
	Slug    string    `json:"slug,omitempty"` // в данной структуре может быть а может и не быть
	Created time.Time `json:"created"`        // время создания ветки
}

type ThreadCreate struct {
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type ThreadPostRequest struct {
	SlugOrId string `json:"slug_or_id"`
	Limit    int32  `json:"limit,omitempty"`
	Since    int64  `json:"since"`
	Sort     string `json:"sort"` // flat, tree или parent_tree
	Desc     bool   `json:"desc"` // флаг сортировки по убыванию
}
