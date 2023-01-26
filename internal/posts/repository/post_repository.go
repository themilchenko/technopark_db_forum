package postRepository

import (
	"fmt"
	"technopark_db_forum/internal/models"
	e "technopark_db_forum/pkg/errors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	CreatePosts(post []models.Post) ([]models.Post, error)
	GetPostByID(id uint64) (models.Post, error)
	UpdatePost(post models.Post) (models.Post, error)

	// GetThreadPosts(id uint64, limit uint64, sort string, since uint64, desc bool) ([]models.Post, error)
	GetThreadPostsFlat(id uint64, limit uint64, since uint64, desk bool) ([]models.Post, error)
	GetThreadPostsTree(id uint64, limit uint64, since uint64, desk bool) ([]models.Post, error)
	GetThreadPostsParentTree(id uint64, limit uint64, since uint64, desk bool) ([]models.Post, error)
}

type Postgres struct {
	DB *sqlx.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

func (p Postgres) CreatePosts(post []models.Post) ([]models.Post, error) {
	var res []models.Post
	query := `INSERT INTO posts (author, created, forum, message, parent, thread) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, author, created, forum, message, parent, thread`
	for _, value := range post {
		var cur models.Post
		err := p.DB.QueryRow(query, value.Author, value.Created, value.Forum, value.Message, value.Parent, value.ThreadID).Scan(&cur.ID, &cur.Author, &cur.Created, &cur.Forum, &cur.Message, &cur.Parent, &cur.ThreadID)
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23503" {
				return make([]models.Post, 0), e.ErrNoAuthorPost

			} else if pgErr.Code == "23505" {
				return make([]models.Post, 0), e.ErrOtherThread
			}
		}
		if err != nil {
			return make([]models.Post, 0), err
		}
		res = append(res, cur)
	}
	return res, nil
}

func (p Postgres) GetPostByID(id uint64) (models.Post, error) {
	query := `SELECT id, author, created, forum, message, parent, thread, is_edited FROM posts WHERE id = $1`
	post := models.Post{}
	err := p.DB.Get(&post, query, id)
	return post, err
}

func (p Postgres) UpdatePost(post models.Post) (models.Post, error) {
	var res models.Post
	query := `UPDATE posts SET message = $1, is_edited=TRUE WHERE id = $2 RETURNING id, author, created, forum, is_edited, message, parent, thread`
	err := p.DB.QueryRow(query, post.Message, post.ID).Scan(&res.ID, &res.Author, &res.Created, &res.Forum, &res.IsEdited, &res.Message, &res.Parent, &res.ThreadID)
	return res, err
}

func (p Postgres) GetThreadPostsFlat(id uint64, limit uint64, since uint64, desk bool) ([]models.Post, error) {
	res := make([]models.Post, 0)
	query := `
		SELECT p.id, p.author, p.created, p.forum, p.is_edited, p.message, p.parent, p.thread FROM posts p WHERE p.thread = $1
	`

	if since != 0 {
		if desk {
			query += fmt.Sprintf(` AND p.id < %d`, since)
		} else {
			query += fmt.Sprintf(` AND p.id > %d`, since)
		}
	}

	if desk {
		query += ` ORDER BY created DESC, p.id DESC`
	} else {
		query += ` ORDER BY created, p.id`
	}

	if limit != 0 {
		query += fmt.Sprintf(` LIMIT %d`, limit)
	}

	err := p.DB.Select(&res, query, id)
	return res, err
}

func (p Postgres) GetThreadPostsTree(id uint64, limit uint64, since uint64, desk bool) ([]models.Post, error) {
	posts := make([]models.Post, 0)

	query := `
		SELECT p.id, p.author, p.created, p.forum, p.is_edited, p.message, p.parent, p.thread
		FROM posts p
		WHERE p.thread = $1
	`

	if since != 0 {
		if desk {
			query += fmt.Sprintf(` AND p.path < (SELECT path FROM posts WHERE id = %d)`, since)
		} else {
			query += fmt.Sprintf(` AND p.path > (SELECT path FROM posts WHERE id = %d)`, since)
		}
	}

	if desk {
		query += ` ORDER BY p.path DESC`
	} else {
		query += ` ORDER BY p.path`
	}

	if limit != 0 {
		query += fmt.Sprintf(` LIMIT %d`, limit)
	}

	err := p.DB.Select(&posts, query, id)
	return posts, err
}

func (p Postgres) GetThreadPostsParentTree(id uint64, limit uint64, since uint64, desk bool) ([]models.Post, error) {
	posts := make([]models.Post, 0)

	query := `SELECT id, author, created, forum, is_edited, message, parent, thread FROM posts`

	if since == 0 {
		if desk {
			query += ` WHERE path[1] IN (SELECT id FROM posts WHERE parent = 0 AND thread = $1 ORDER BY id DESC LIMIT $2)`
		} else {
			query += ` WHERE path[1] IN (SELECT id FROM posts WHERE parent = 0 AND thread = $1 ORDER BY id ASC LIMIT $2)`
		}
	} else {
		if desk {
			query += ` WHERE path[1] IN (SELECT id FROM posts WHERE parent = 0 AND thread = $1 AND id < (SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3)`
		} else {
			query += ` WHERE path[1] IN (SELECT id FROM posts WHERE parent = 0 AND thread = $1 AND id > (SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3)`
		}
	}

	if desk {
		query += ` ORDER BY path[1] DESC, path`
	} else {
		query += ` ORDER BY path[1] ASC, path`
	}

	if since != 0 {
		err := p.DB.Select(&posts, query, id, since, limit)
		return posts, err
	}
	err := p.DB.Select(&posts, query, id, limit)
	return posts, err
}
