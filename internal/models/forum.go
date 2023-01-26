package models

type Forum struct {
	Slug         string `json:"slug" db:"slug"`
	Title        string `json:"title" db:"title"`
	UserNickname string `json:"user" db:"user_nick"`

	PostsCount   uint64 `json:"posts" db:"posts"`
	ThreadsCount uint64 `json:"threads" db:"threads"`
}

type ForumUser struct {
	ForumSlug    string `json:"forum" db:"forum"`
	UserNickname string `json:"user" db:"user_nick"`
}
