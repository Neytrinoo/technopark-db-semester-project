package models

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type UserUpdate struct {
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UserCreate struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
