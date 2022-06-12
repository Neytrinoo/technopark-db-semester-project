package models

type Vote struct {
	Nickname string `json:"nickname"` // автор голоса
	Thread   int64  `json:"thread"`
	Voice    int32  `json:"voice"` // -1 или 1, голос
}

type VoteCreate struct {
	Nickname string `json:"nickname"` // автор голоса
	Voice    int32  `json:"voice"`    // -1 или 1, голос
}
