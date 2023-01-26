package usecase

import (
	"technopark_db_forum/internal/models"
	serviceRepository "technopark_db_forum/internal/service/repository"
)

type ServiceUsecase interface {
	GetStatus() (models.ServiceStatus, error)
	Clear() error
}

type usecase struct {
	serviceRepository serviceRepository.ServiceRepository
}

func NewServiceUsecase(serviceRepo serviceRepository.ServiceRepository) ServiceUsecase {
	return &usecase{
		serviceRepository: serviceRepo,
	}
}

func (u usecase) GetStatus() (models.ServiceStatus, error) {
	status, err := u.serviceRepository.GetStatus()
	if err != nil {
		return models.ServiceStatus{}, err
	}
	return status, nil
}

func (u usecase) Clear() error {
	err := u.serviceRepository.Clear()
	if err != nil {
		return err
	}
	return nil
}
