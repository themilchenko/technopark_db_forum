package forumRepository

import (
	"technopark_db_forum/internal/models"
	"technopark_db_forum/pkg/errors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type ForumRepository interface {
	CreateForum(forum models.Forum) (models.Forum, error)
	CreateForumUser(forum, user string) (models.ForumUser, error)
	GetForumBySlug(slug string) (models.Forum, error)
	GetForumUsers(slug string, options models.ThreadOptions) ([]models.User, error)
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

func (p *Postgres) CreateForum(forum models.Forum) (models.Forum, error) {
	var res models.Forum
	query := `INSERT INTO forums (slug, title, user_nick) VALUES ($1, $2, $3) RETURNING slug, title, user_nick, posts, threads`
	err := p.DB.QueryRow(query, forum.Slug, forum.Title, forum.UserNickname).Scan(&res.Slug, &res.Title, &res.UserNickname, &res.PostsCount, &res.ThreadsCount)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return models.Forum{}, errors.ErrUserNotFound
		}

		return models.Forum{}, err
	}
	return res, err
}

func (p *Postgres) CreateForumUser(forum, user string) (models.ForumUser, error) {
	var res models.ForumUser
	query := `INSERT INTO forum_users (forum, user_nick) VALUES ($1, $2) RETURNING forum, user_nick`
	err := p.DB.QueryRow(query, forum, user).Scan(&res.ForumSlug, &res.UserNickname)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				return res, nil
			}
		}
		return models.ForumUser{}, err
	}
	return res, err
}

func (p *Postgres) GetForumBySlug(slug string) (models.Forum, error) {
	query := `SELECT slug, title, user_nick, posts, threads FROM forums WHERE slug = $1`
	forum := models.Forum{}
	err := p.DB.Get(&forum, query, slug)
	return forum, err
}

func (p *Postgres) GetUsers(slug string) ([]models.User, error) {
	query := `SELECT u.nickname, u.fullname, u.email, u.about FROM forum_users
		JOIN users u ON forum_users.user_nick = u.nickname
		WHERE forum_users.forum_slug = $1`
	users := []models.User{}
	err := p.DB.Select(&users, query, slug)
	return users, err
}

func (p *Postgres) GetForumUsers(slug string, options models.ThreadOptions) ([]models.User, error) {
	query := `SELECT users.nickname, users.fullname, users.email, users.about FROM users JOIN forum_users ON forum_users.user_nick = users.nickname WHERE forum_users.forum = $1`
	if options.Since != "" && options.Desc {
		query += ` AND users.nickname < $2`
	} else if options.Since != "" && !options.Desc {
		query += ` AND users.nickname > $2`
	}
	query += ` ORDER BY users.nickname`
	if options.Desc {
		query += ` DESC`
	}
	if options.Limit != 0 {
		if options.Since != "" {
			query += ` LIMIT $3`
		} else {
			query += ` LIMIT $2`
		}
	}

	users := make([]models.User, 0)
	var err error
	if options.Limit != 0 && options.Since != "" {
		err := p.DB.Select(&users, query, slug, options.Since, options.Limit)
		if err != nil {
			return nil, err
		}
	} else if options.Limit != 0 && options.Since == "" {
		err := p.DB.Select(&users, query, slug, options.Limit)
		if err != nil {
			return nil, err
		}
	}

	return users, err
}
