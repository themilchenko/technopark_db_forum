package usecase

import (
	"database/sql"
	"strconv"
	forumRepository "technopark_db_forum/internal/forum/repository"
	"technopark_db_forum/internal/models"
	"technopark_db_forum/internal/posts/repository"
	"technopark_db_forum/internal/thread/repository"
	"technopark_db_forum/internal/users/repository"
	e "technopark_db_forum/pkg/errors"
	"time"
)

type PostUsecase interface {
	CreatePosts(posts []models.Post, slugOrID string) ([]models.Post, error)
	GetPostByID(id uint64) (models.Post, error)
	GetPostByIDRelared(id uint64, related []string) (models.PostFull, error)
	UpdatePost(post models.Post) (models.Post, error)
	GetThreadPosts(slugOrID string, limit uint64, sort string, since uint64, desk bool) ([]models.Post, error)
}

type usecase struct {
	postRepository   postRepository.PostRepository
	userRepository   userRepository.UserRepository
	threadRepository threadRepository.ThreadRepository
	forumRepository  forumRepository.ForumRepository
}

func NewPostUsecase(postRepo postRepository.PostRepository, userRepo userRepository.UserRepository, threadRepo threadRepository.ThreadRepository, forumRepo forumRepository.ForumRepository) PostUsecase {
	return &usecase{
		postRepository:   postRepo,
		userRepository:   userRepo,
		threadRepository: threadRepo,
		forumRepository:  forumRepo,
	}
}

func (u usecase) CreatePosts(posts []models.Post, slugOrID string) ([]models.Post, error) {
	try, err := strconv.ParseUint(slugOrID, 10, 64)
	var thread models.Thread
	if err != nil {
		thread, err = u.threadRepository.GetThreadBySlug(slugOrID)
		if err != nil {
			return nil, e.ErrNoRowsSlug
		}
	} else {
		thread, err = u.threadRepository.GetThreadByID(try)
		if err != nil {
			return nil, e.ErrNoRowsID
		}
	}


	curTime := time.Now().Format(time.RFC3339)
	for index := range posts {
		_, err := u.userRepository.GetUserByNickname(posts[index].Author)
		if err != nil {
			return nil, e.ErrNoAuthorPost
		}

		if posts[index].Parent != 0 {
			parent, err := u.postRepository.GetPostByID(posts[index].Parent)

			
			if posts[index].ThreadID != thread.ID {
				return nil, e.ErrConflict
			}
			if err == sql.ErrNoRows {
				return nil, e.ErrOtherThread
			}
			if parent.Forum != thread.Forum {
				return nil, e.ErrOtherThread
			}
		}

		posts[index].ThreadID = thread.ID
		posts[index].Forum = thread.Forum
		posts[index].Created = curTime
	}

	res, err := u.postRepository.CreatePosts(posts)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return make([]models.Post, 0), nil
	}

	for index := range res {
		_, err := u.forumRepository.CreateForumUser(res[index].Forum, res[index].Author)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (u usecase) GetPostByIDRelared(id uint64, related []string) (models.PostFull, error) {
	post, err := u.postRepository.GetPostByID(id)
	if err != nil {
		return models.PostFull{}, err
	}

	var postFull models.PostFull
	postFull.Post = &post

	for _, rel := range related {
		switch rel {
		case "user":
			user, err := u.userRepository.GetUserByNickname(post.Author)
			if err != nil {
				return models.PostFull{}, err
			}
			postFull.Author = &user
		case "forum":
			forum, err := u.forumRepository.GetForumBySlug(post.Forum)
			if err != nil {
				return models.PostFull{}, err
			}
			postFull.Forum = &forum
		case "thread":
			thread, err := u.threadRepository.GetThreadByID(post.ThreadID)
			if err != nil {
				return models.PostFull{}, err
			}
			postFull.Thread = &thread
		}
	}

	return postFull, nil
}

func (u usecase) GetPostByID(id uint64) (models.Post, error) {
	post, err := u.postRepository.GetPostByID(id)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (u usecase) UpdatePost(post models.Post) (models.Post, error) {
	if post.Message == "" {
		return u.postRepository.GetPostByID(post.ID)
	}
	res, err := u.postRepository.UpdatePost(post)
	if err != nil {
		return models.Post{}, err
	}
	return res, nil
}

func (u usecase) GetThreadPosts(slugOrID string, limit uint64, sort string, since uint64, desk bool) ([]models.Post, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	var th models.Thread
	if err != nil {
		th, err = u.threadRepository.GetThreadBySlug(slugOrID)
		if err != nil {
			return nil, e.ErrNoRowsSlug
		}
	} else {
		th, err = u.threadRepository.GetThreadByID(id)
		if err != nil {
			return nil, e.ErrNoRowsID
		}
	}

	switch sort {
	case "flat":
		res, err := u.postRepository.GetThreadPostsFlat(th.ID, limit, since, desk)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "tree":
		res, err := u.postRepository.GetThreadPostsTree(th.ID, limit, since, desk)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "parent_tree":
		res, err := u.postRepository.GetThreadPostsParentTree(th.ID, limit, since, desk)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		res, err := u.postRepository.GetThreadPostsFlat(th.ID, limit, since, desk)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
