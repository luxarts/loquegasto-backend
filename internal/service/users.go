package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type UsersService interface {
	Create(userDTO *domain.UserDTO) (*domain.UserDTO, error)
	GetByID(id int64) (*domain.UserDTO, error)
	Update(userDTO *domain.UserDTO) (*domain.UserDTO, error)
	Delete(id int64) error
}
type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersService{
		repo: repo,
	}
}
func (s *usersService) Create(userDTO *domain.UserDTO) (*domain.UserDTO, error) {
	user := userDTO.ToUser()

	user, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user.ToDTO(), nil
}
func (s *usersService) GetByID(id int64) (*domain.UserDTO, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}
func (s *usersService) Update(userDTO *domain.UserDTO) (*domain.UserDTO, error) {
	user := userDTO.ToUser()

	user, err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user.ToDTO(), nil
}
func (s *usersService) Delete(id int64) error {
	return s.repo.Delete(id)
}
