package models

type ServiceStatus struct {
	UsersCount  uint64 `json:"user"`
	ForumsCount uint64 `json:"forum"`
	ThreadsCount uint64 `json:"thread"`
	PostsCount uint64 `json:"post"`
}