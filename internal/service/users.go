package service

import (
	"github.com/google/uuid"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/utils/jwt"
	"time"
)

type UsersService interface {
	Create(req *domain.UserCreateRequest) (*domain.UserCreateResponse, error)
	AuthWithTelegram(req *domain.UserAuthWithTelegramRequest) (*domain.UserAuthWithTelegramResponse, error)
}
type usersService struct {
	repo repository.UsersRepository
}

func NewUsersService(repo repository.UsersRepository) UsersService {
	return &usersService{
		repo: repo,
	}
}
func (s *usersService) Create(req *domain.UserCreateRequest) (*domain.UserCreateResponse, error) {
	user := req.ToUser()

	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()

	user, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}
func (s *usersService) AuthWithTelegram(req *domain.UserAuthWithTelegramRequest) (*domain.UserAuthWithTelegramResponse, error) {
	u, err := s.repo.GetByChatID(req.ChatID)
	if err != nil {
		return nil, err
	}

	token := jwt.GenerateToken(nil, &jwt.Payload{Subject: u.ID})

	return &domain.UserAuthWithTelegramResponse{
		AccessToken: token,
	}, nil
}
