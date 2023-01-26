package userRepository

import (
	"technopark_db_forum/internal/models"
	e "technopark_db_forum/pkg/errors"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByNickname(nickname string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	GetUsersByEmailNickname(email, nickname string) ([]models.User, error)
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

func (p Postgres) GetUsersByEmailNickname(email, nickname string) ([]models.User, error) {
	query := `SELECT nickname, fullname, email, about FROM users WHERE email = $1 OR nickname = $2`
	users := []models.User{}
	err := p.DB.Select(&users, query, email, nickname)
	return users, err
}

func (p Postgres) CreateUser(user models.User) (models.User, error) {
	var res models.User
	query := `INSERT INTO users (nickname, fullname, email, about) VALUES ($1, $2, $3, $4) RETURNING *`
	err := p.DB.QueryRowx(query, user.Nickname, user.FullName, user.Email, user.About).Scan(&res.Nickname, &res.FullName, &res.Email, &res.About)
	return res, err
}

func (p Postgres) GetUserByNickname(nickname string) (models.User, error) {
	query := `SELECT nickname, fullname, email, about FROM users WHERE nickname = $1`
	user := models.User{}
	err := p.DB.Get(&user, query, nickname)
	return user, err
}

func (p Postgres) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT nickname, fullname, email, about FROM users WHERE email = $1`
	user := models.User{}
	err := p.DB.Get(&user, query, email)
	return user, err
}

func (p Postgres) UpdateUser(user models.User) (models.User, error) {
	var res models.User
	query := `UPDATE users SET fullname = $1, email = $2, about = $3 WHERE nickname = $4 RETURNING fullname, email, about, nickname`
	err := p.DB.QueryRow(query, user.FullName, user.Email, user.About, user.Nickname).Scan(&res.FullName, &res.Email, &res.About, &res.Nickname)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code == "23505" {
				if pgErr.Constraint == "users_email_key" {
					return models.User{}, e.ErrConflictEmail
				} else {
					return models.User{}, e.ErrConflictNickname
				}
			}
		}
		return models.User{}, err
	}
	return res, err
}
