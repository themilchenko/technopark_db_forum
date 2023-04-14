package models

type Post struct {
	ID        uint64    `json:"id,omitempty" db:"id"`
	Author    string    `json:"author,omitempty" db:"author"`
	Forum     string    `json:"forum,omitempty" db:"forum"`
	ThreadID  uint64    `json:"thread,omitempty" db:"thread"`
	Message   string    `json:"message,omitempty" db:"message"`
	IsEdited  bool      `json:"isEdited,omitempty" db:"is_edited"`
	Created   string    `json:"created,omitempty" db:"created"`
	Parent    uint64    `json:"parent,omitempty" db:"parent"`
	Path      []uint64  `json:"path,omitempty" db:"path"`
	TreeLevel uint64    `json:"tree_level,omitempty" db:"tree_level"`
}

type PostFull struct {
	Post   *Post   `json:"post,omitempty"`
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}
