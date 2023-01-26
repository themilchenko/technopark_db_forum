package usecase

import (
	"database/sql"
	"strconv"
	forumRepository "technopark_db_forum/internal/forum/repository"
	"technopark_db_forum/internal/models"
	threadRepository "technopark_db_forum/internal/thread/repository"
	userRepository "technopark_db_forum/internal/users/repository"
	e "technopark_db_forum/pkg/errors"
	"time"

	"github.com/jinzhu/copier"
)

type ThreadUsecase interface {
	CreateThread(thread models.Thread) (models.Thread, error)
	GetThreadBySlug(slugOrID string) (models.Thread, error)
	GetThreadMsgsBySlug(slug string, since time.Time, options models.ThreadOptions) ([]models.Thread, error)
	GetThreadByID(id uint64) (models.Thread, error)
	UpdateThread(thread models.Thread, slugOrID string) (models.ThreadNoVotes, error)
	CreateVote(vote models.Vote, slugOrID string) (models.Thread, error)
	GetThread(slugOrID string) (models.Thread, error)
}

type usecase struct {
	threadRepository threadRepository.ThreadRepository
	userRepository   userRepository.UserRepository
	forumRepository  forumRepository.ForumRepository
}

func NewThreadUsecase(threadRepo threadRepository.ThreadRepository, userRepo userRepository.UserRepository, forumRepo forumRepository.ForumRepository) ThreadUsecase {
	return &usecase{
		threadRepository: threadRepo,
		userRepository:   userRepo,
		forumRepository:  forumRepo,
	}
}

func (u usecase) CreateThread(thread models.Thread) (models.Thread, error) {
	forum, err := u.forumRepository.GetForumBySlug(thread.Forum)
	if err != nil {
		return models.Thread{}, e.ErrThreadNotFound
	}

	_, err = u.userRepository.GetUserByNickname(thread.Author)
	if err != nil {
		return models.Thread{}, e.ErrConflictNickname
	}

	if thread.Slug != "" {
		th, err := u.threadRepository.GetThreadBySlug(thread.Slug)
		if err == nil {
			return th, e.ErrDuplicate
		}
		if err != nil && err != sql.ErrNoRows {
			return models.Thread{}, err
		}
	}

	thread.Forum = forum.Slug

	res, err := u.threadRepository.CreateThread(thread)
	if err != nil {
		return models.Thread{}, err
	}

	_, err = u.forumRepository.CreateForumUser(thread.Forum, thread.Author)
	if err != nil {
		return models.Thread{}, err
	}

	return res, nil
}

func (u usecase) GetThreadBySlug(slugOrID string) (models.Thread, error) {
	thread, err := u.threadRepository.GetThreadBySlug(slugOrID)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}

func (u usecase) GetThreadMsgsBySlug(slugOrID string, since time.Time, options models.ThreadOptions) ([]models.Thread, error) {
	_, err := u.forumRepository.GetForumBySlug(slugOrID)
	if err != nil {
		return nil, err
	}

	threads, err := u.threadRepository.GetThreadMsgs(slugOrID, since, options)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (u usecase) GetThreadByID(id uint64) (models.Thread, error) {
	thread, err := u.threadRepository.GetThreadByID(id)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}

func (u usecase) UpdateThread(thread models.Thread, slugOrID string) (models.ThreadNoVotes, error) {
	try, err := strconv.ParseUint(slugOrID, 10, 64)
	var th models.Thread
	if err == nil {
		th, err = u.threadRepository.GetThreadByID(try)
		if err != nil {
			return models.ThreadNoVotes{}, e.ErrNoRowsID
		}
	} else {
		th, err = u.threadRepository.GetThreadBySlug(slugOrID)
		if err != nil {
			return models.ThreadNoVotes{}, e.ErrNoRowsSlug
		}
		thread.Slug = slugOrID
	}

	if thread.Message == "" && thread.Title == "" {
		return models.ThreadNoVotes{
			ID:      th.ID,
			Title:   th.Title,
			Author:  th.Author,
			Forum:   th.Forum,
			Message: th.Message,
			Slug:    th.Slug,
			Created: th.Created.UTC(),
		}, nil
	}

	if err = copier.CopyWithOption(&th, &thread, copier.Option{IgnoreEmpty: true}); err != nil {
		return models.ThreadNoVotes{}, err
	}

	thread.ID = th.ID
	thread.Slug = th.Slug

	res, err := u.threadRepository.UpdateThread(th)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ThreadNoVotes{
				ID:      th.ID,
				Title:   th.Title,
				Author:  th.Author,
				Forum:   th.Forum,
				Message: th.Message,
				Slug:    th.Slug,
				Created: th.Created.UTC(),
			}, nil
		}
		return models.ThreadNoVotes{}, err
	}

	return res, nil
}

func (u usecase) CreateVote(vote models.Vote, slugOrID string) (models.Thread, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err != nil {
		return u.threadRepository.VoteBySlug(slugOrID, vote)
	}

	return u.threadRepository.VoteByID(id, vote)
}

func (u usecase) GetThread(slugOrID string) (models.Thread, error) {
	try, err := strconv.ParseUint(slugOrID, 10, 64)
	if err == nil {
		thread, err := u.threadRepository.GetThreadByID(try)
		if err != nil {
			return models.Thread{}, err
		}
		return thread, nil
	}

	thread, err := u.threadRepository.GetThreadBySlug(slugOrID)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}
