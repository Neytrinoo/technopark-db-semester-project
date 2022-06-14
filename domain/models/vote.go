package models

type Vote struct {
	Nickname string `json:"nickname" db:"nickname"` // автор голоса
	Thread   int64  `json:"thread" db:"thread"`
	Voice    int32  `json:"voice" db:"voice"` // -1 или 1, голос
}

type VoteCreate struct {
	Nickname string `json:"nickname" db:"nickname"` // автор голоса
	Voice    int32  `json:"voice" db:"voice"`       // -1 или 1, голос
}
