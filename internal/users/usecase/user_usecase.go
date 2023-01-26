package usecase

import (
	"database/sql"
	"errors"
	"technopark_db_forum/internal/models"
	userRepository "technopark_db_forum/internal/users/repository"
	e "technopark_db_forum/pkg/errors"

	"github.com/jinzhu/copier"
)

type UsersUsecase interface {
	CreateUser(user models.User) ([]models.User, error)
	GetUsersByEmailNickname(email, nickname string) ([]models.User, error)
	GetUserByNickname(nickname string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
}

type usecase struct {
	userRepository userRepository.UserRepository
}

func NewUserUsecase(userRepo userRepository.UserRepository) UsersUsecase {
	return &usecase{
		userRepository: userRepo,
	}
}

func (u usecase) GetUsersByEmailNickname(email, nickname string) ([]models.User, error) {
	users, err := u.userRepository.GetUsersByEmailNickname(email, nickname)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (u usecase) CreateUser(user models.User) ([]models.User, error) {
	users, err := u.userRepository.GetUsersByEmailNickname(user.Email, user.Nickname)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return []models.User{}, err
		}
	}
	if len(users) != 0 {
		return users, e.ErrDuplicate
	}

	res, err := u.userRepository.CreateUser(user)
	if err != nil {
		return []models.User{}, err
	}
	return []models.User{res}, nil
}

func (u usecase) GetUserByNickname(nickname string) (models.User, error) {
	user, err := u.userRepository.GetUserByNickname(nickname)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u usecase) GetUserByEmail(email string) (models.User, error) {
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u usecase) UpdateUser(user models.User) (models.User, error) {
	us, err := u.userRepository.GetUserByNickname(user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	if err = copier.CopyWithOption(&us, &user, copier.Option{IgnoreEmpty: true}); err != nil {
		return models.User{}, err
	}
	
	res, err := u.userRepository.UpdateUser(us)
	if err != nil {
		return models.User{}, err
	}
	return res, nil
}
