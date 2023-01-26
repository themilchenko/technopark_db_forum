package models

type User struct {
	Nickname string `json:"nickname" db:"nickname"`
	FullName string `json:"fullname" db:"fullname"`
	Email    string `json:"email" db:"email"`
	About    string `json:"about" db:"about"`
}
