package errors

import "errors"

var (
	ErrDuplicate        = errors.New("duplicate")
	ErrNoParent         = errors.New("any post has no parent")
	ErrConflict         = errors.New("user conflict")
	ErrConflictEmail    = errors.New("user conflict")
	ErrConflictNickname = errors.New("user conflict")
	ErrUserNotFound     = errors.New("user not found")
	ErrForumNotFound    = errors.New("forum not found")
	ErrThreadNotFound   = errors.New("thread not found")
	ErrPostNotFound     = errors.New("post not found")
	ErrNotFound         = errors.New("not found")
	ErrInternal         = errors.New("internal error")
	ErrOtherThread      = errors.New("other thread")
	ErrNoAuthorPost     = errors.New("no author post")
	ErrNoRowsID         = errors.New("no rows in result set")
	ErrNoRowsSlug       = errors.New("no rows in result set")
	ErrNoRows = errors.New("no rows in result set")
)
