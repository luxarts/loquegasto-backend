package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type UsersService interface {
	Create(userDTO *domain.UserDTO) (*domain.UserDTO, error)
	GetByID(id int) (*domain.UserDTO, error)
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
func (s *usersService) GetByID(id int) (*domain.UserDTO, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}
