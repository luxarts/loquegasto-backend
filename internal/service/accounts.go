package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type AccountsService interface {
	Create(accountDTO *domain.AccountDTO) (*domain.AccountDTO, error)
	GetByName(userID int, name string) (*domain.AccountDTO, error)
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
func (s *accountsService) GetByName(userID int, name string) (*domain.AccountDTO, error) {
	accounts, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	name = sanitizeName(name)

	for _, a := range *accounts {
		if sanitizeName(a.Name) == name {
			return a.ToDTO(), nil
		}
	}

	return nil, nil
}

// Utils
func sanitizeName(s string) string {
	// Remove leading/trailing spaces
	s = strings.TrimSpace(s)
	// Converts letters with tildes to normal
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)
	// Convert string to lowercase
	s = strings.ToLower(s)

	return s
}
