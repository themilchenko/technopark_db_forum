package serviceRepository

import (
	_ "github.com/lib/pq"
	"technopark_db_forum/internal/models"

	"github.com/jmoiron/sqlx"
)

type ServiceRepository interface {
	GetStatus() (models.ServiceStatus, error)
	Clear() error
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

 func (p Postgres) GetStatus() (models.ServiceStatus, error) {
	var res models.ServiceStatus
	
	query := `SELECT COUNT(*) FROM users`
	err := p.DB.Get(&res.UsersCount, query)
	if err != nil {
		return models.ServiceStatus{}, err
	}

	query = `SELECT COUNT(*) FROM forums`
	err = p.DB.Get(&res.ForumsCount, query)
	if err != nil {
		return models.ServiceStatus{}, err
	}

	query = `SELECT COUNT(*) FROM threads`
	err = p.DB.Get(&res.ThreadsCount, query)
	if err != nil {
		return models.ServiceStatus{}, err
	}

	query = `SELECT COUNT(*) FROM posts`
	err = p.DB.Get(&res.PostsCount, query)
	if err != nil {
		return models.ServiceStatus{}, err
	}

	return res, nil
}

func (p Postgres) Clear() error {
	_, err := p.DB.Exec(`TRUNCATE forums CASCADE;
						 TRUNCATE threads CASCADE;
						 TRUNCATE posts CASCADE;
						 TRUNCATE votes CASCADE;
						 TRUNCATE forum_users CASCADE;
						 TRUNCATE users CASCADE;`)
	return err
}
