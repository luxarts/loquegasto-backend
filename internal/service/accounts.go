package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type AccountsService interface {
	Create(accountDTO *domain.AccountDTO) (*domain.AccountDTO, error)
}
type accountsService struct {
	repo repository.AccountRepository
}

func NewAccountsService(repo repository.AccountRepository) AccountsService {
	return &accountsService{
		repo: repo,
	}
}
func (s *accountsService) Create(accountDTO *domain.AccountDTO) (*domain.AccountDTO, error) {
	account := accountDTO.ToAccount()

	account, err := s.repo.Create(account)
	if err != nil {
		return nil, err
	}

	return account.ToDTO(), nil
}
