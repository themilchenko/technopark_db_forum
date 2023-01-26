package usecase

import (
	"database/sql"
	"technopark_db_forum/internal/forum/repository"
	"technopark_db_forum/internal/models"
	"technopark_db_forum/internal/users/repository"
	e "technopark_db_forum/pkg/errors"
)

type ForumUsecase interface {
	CreateForum(forum models.Forum) (models.Forum, error)
	GetForumBySlug(slug string) (models.Forum, error)
	GetForumUsersBySlug(slug string, options models.ThreadOptions) ([]models.User, error)
}

type usecase struct {
	forumRepository forumRepository.ForumRepository
	userRepository  userRepository.UserRepository
}

func NewUserUsecase(forumRepo forumRepository.ForumRepository, userRepo userRepository.UserRepository) ForumUsecase {
	return &usecase{
		forumRepository: forumRepo,
		userRepository:  userRepo,
	}
}

func (u usecase) CreateForum(forum models.Forum) (models.Forum, error) {
	user, err := u.userRepository.GetUserByNickname(forum.UserNickname)
	if err != nil {
		return models.Forum{}, err
	}

	f, err := u.forumRepository.GetForumBySlug(forum.Slug)
	if err != nil && err != sql.ErrNoRows {
		return models.Forum{}, err
	}
	if err == nil {
		return f, e.ErrDuplicate
	}

	forum.UserNickname = user.Nickname
	res, err := u.forumRepository.CreateForum(forum)
	if err != nil {
		if err == e.ErrUserNotFound {
			return models.Forum{}, e.ErrUserNotFound
		}
		return models.Forum{}, err
	}
	return res, nil
}

func (u usecase) GetForumBySlug(slug string) (models.Forum, error) {
	forum, err := u.forumRepository.GetForumBySlug(slug)
	if err != nil {
		return models.Forum{}, err
	}
	return forum, nil
}

func (u usecase) GetForumUsersBySlug(slug string, options models.ThreadOptions) ([]models.User, error) {
	_, err := u.forumRepository.GetForumBySlug(slug)
	if err != nil {
		return []models.User{}, err
	}

	users, err := u.forumRepository.GetForumUsers(slug, options)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}
