package models

import "time"

type Thread struct {
	ID      uint64    `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Author  string    `json:"author" db:"author"`
	Forum   string    `json:"forum" db:"forum"`
	Message string    `json:"message" db:"message"`
	Votes   int64     `json:"votes" db:"votes"`
	Slug    string    `json:"slug" db:"slug"`
	Created time.Time `json:"created" db:"created"`
}

type ThreadNoVotes struct {
	ID      uint64 `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Author  string `json:"author" db:"author"`
	Forum   string `json:"forum" db:"forum"`
	Message string `json:"message" db:"message"`
	Slug    string `json:"slug" db:"slug"`
	Created time.Time `json:"created" db:"created"`
}

type Vote struct {
	ID         uint64 `json:"id" db:"id"`
	Nickname   string `json:"nickname" db:"nickname"`
	VoiceValue int64  `json:"voice" db:"voice"`
	ThreadID   uint64 `json:"thread" db:"thread"`
}

type ThreadOptions struct {
	Limit  uint64
	Since  string
	Desc   bool
	SortBy string
}
