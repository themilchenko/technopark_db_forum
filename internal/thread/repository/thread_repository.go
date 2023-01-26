package threadRepository

import (
	"technopark_db_forum/internal/models"
	e "technopark_db_forum/pkg/errors"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type ThreadRepository interface {
	CreateThread(thread models.Thread) (models.Thread, error)
	GetThreadBySlug(slug string) (models.Thread, error)
	GetThreadByID(id uint64) (models.Thread, error)
	GetThreadMsgs(slugOrID string, since time.Time, options models.ThreadOptions) ([]models.Thread, error)
	UpdateThread(thread models.Thread) (models.ThreadNoVotes, error)

	VoteBySlug(slug string, v models.Vote) (models.Thread, error)
	VoteByID(id uint64, v models.Vote) (models.Thread, error)
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

func (p Postgres) CreateThread(thread models.Thread) (models.Thread, error) {
	var res models.Thread
	query := `INSERT INTO threads (slug, author, forum, title, message, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, slug, author, forum, title, message, votes, created`
	err := p.DB.QueryRowx(query, thread.Slug, thread.Author, thread.Forum, thread.Title, thread.Message, thread.Created).Scan(&res.ID, &res.Slug, &res.Author, &res.Forum, &res.Title, &res.Message, &res.Votes, &res.Created)
	return res, err
}

func (p Postgres) InsertThread(thread models.Thread) error {
	query := `INSERT INTO threads (slug, author, forum, title, message, created) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := p.DB.Exec(query, thread.Slug, thread.Author, thread.Forum, thread.Title, thread.Message, thread.Created)
	return err
}

func (p Postgres) GetThreadBySlug(slug string) (models.Thread, error) {
	query := `SELECT id, slug, author, forum, title, message, votes, created FROM threads WHERE slug = $1`
	thread := models.Thread{}
	err := p.DB.Get(&thread, query, slug)
	return thread, err
}

func (p Postgres) GetThreadByID(id uint64) (models.Thread, error) {
	query := `SELECT id, slug, author, forum, title, message, votes, created FROM threads WHERE id = $1`
	thread := models.Thread{}
	err := p.DB.Get(&thread, query, id)
	return thread, err
}

func (p Postgres) GetThreadMsgs(slugOrID string, since time.Time, options models.ThreadOptions) ([]models.Thread, error) {
	query := `SELECT threads.id, threads.slug, threads.author, threads.forum, threads.title, threads.message, threads.votes, threads.created FROM threads WHERE threads.forum = $1`
	isTime := false
	if since != (time.Time{}) {
		isTime = true
		if options.Desc {
			query += ` AND created <= $2`
		} else {
			query += ` AND created >= $2`
		}
	}
	query += ` ORDER BY created`
	if options.Desc {
		query += ` DESC`
	}
	if isTime {
		query += ` LIMIT $3`
	} else {
		query += ` LIMIT $2`
	}

	threads := make([]models.Thread, 0)
	var err error
	if isTime {
		err = p.DB.Select(&threads, query, slugOrID, since, options.Limit)
	} else {
		err = p.DB.Select(&threads, query, slugOrID, options.Limit)
	}
	return threads, err
}

func (p Postgres) UpdateThread(thread models.Thread) (models.ThreadNoVotes, error) {
	var res models.ThreadNoVotes
	query := `UPDATE threads SET title = $1, message = $2 WHERE slug = $3 RETURNING id, slug, author, forum, title, message, created`
	err := p.DB.QueryRowx(query, thread.Title, thread.Message, thread.Slug).Scan(&res.ID, &res.Slug, &res.Author, &res.Forum, &res.Title, &res.Message, &res.Created)
	return res, err
}

func (p Postgres) VoteBySlug(slug string, v models.Vote) (models.Thread, error) {
	thread := models.Thread{
		Slug: slug,
	}

	err := p.DB.Get(
		&thread,
		`
			SELECT id, author, created, forum, message, slug, title, votes
			FROM threads
			WHERE slug = $1
		`,
		slug,
	)
	if err != nil {
		return thread, err
	}

	_, err = p.DB.Exec(
		`
			INSERT INTO votes (nickname, thread, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (nickname, thread) DO UPDATE
			SET voice = $3
		`,
		v.Nickname,
		thread.ID,
		v.VoiceValue,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return models.Thread{}, e.ErrUserNotFound
		}
		return models.Thread{}, err
	}

	err = p.DB.Get(
		&thread,
		`
			SELECT id, author, created, forum, message, slug, title, votes
			FROM threads
			WHERE slug = $1
		`,
		slug,
	)
	if err != nil {
		return thread, err
	}

	return thread, nil
}

func (p *Postgres) VoteByID(id uint64, v models.Vote) (models.Thread, error) {
	thread := models.Thread{
		ID: id,
	}

	err := p.DB.Get(
		&thread,
		`
			SELECT id, author, created, forum, message, slug, title, votes
			FROM threads
			WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return thread, err
	}

	_, err = p.DB.Exec(
		`
			INSERT INTO votes (nickname, thread, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (nickname, thread) DO UPDATE
			SET voice = $3
		`,
		v.Nickname,
		thread.ID,
		v.VoiceValue,
	)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23503" {
			return models.Thread{}, e.ErrUserNotFound
		}
		return models.Thread{}, err
	}

	err = p.DB.Get(
		&thread,
		`
			SELECT id, author, created, forum, message, slug, title, votes
			FROM threads
			WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return thread, err
	}

	return thread, nil
}
