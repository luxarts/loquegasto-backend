package service

import (
	"github.com/google/uuid"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"time"
)

type UsersService interface {
	Create(userDTO *domain.UserCreateRequest) (*domain.UserCreateResponse, error)
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
