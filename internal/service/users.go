package service

import (
	"github.com/google/uuid"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"time"
)

type UsersService interface {
	Create(userDTO *domain.UserCreateRequest) (*domain.UserCreateResponse, error)
	GetByID(id int64) (*domain.UserCreateResponse, error)
	Update(userDTO *domain.UserCreateRequest) (*domain.UserCreateResponse, error)
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
func (s *usersService) Create(userReq *domain.UserCreateRequest) (*domain.UserCreateResponse, error) {
	user := userReq.ToUser()

	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()

	user, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}
func (s *usersService) GetByID(id int64) (*domain.UserCreateResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}
func (s *usersService) Update(userDTO *domain.UserCreateRequest) (*domain.UserCreateResponse, error) {
	user := userDTO.ToUser()

	user, err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}
func (s *usersService) Delete(id int64) error {
	return s.repo.Delete(id)
}
