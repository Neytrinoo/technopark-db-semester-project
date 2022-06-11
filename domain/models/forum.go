package models

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
