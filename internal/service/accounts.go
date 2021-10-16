package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"net/http"
	"strings"
	"unicode"

	"github.com/luxarts/jsend-go"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type AccountsService interface {
	Create(accountDTO *domain.AccountDTO) (*domain.AccountDTO, error)
	GetByName(userID int, name string) (*domain.AccountDTO, error)
	GetByID(userID int, id int) (*domain.AccountDTO, error)
	GetAll(userID int) (*[]domain.AccountDTO, error)
	UpdateByID(accountDTO *domain.AccountDTO) (*domain.AccountDTO, error)
	Delete(id int, userID int) error
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
func (s *accountsService) GetByID(userID int, id int) (*domain.AccountDTO, error) {
	account, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, jsend.NewError("forbidden", nil, http.StatusForbidden)
	}

	return account.ToDTO(), nil
}
func (s *accountsService) GetAll(userID int) (*[]domain.AccountDTO, error) {
	accounts, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var accountsDTOs []domain.AccountDTO
	for _, a := range *accounts {
		accountsDTOs = append(accountsDTOs, *a.ToDTO())
	}

	return &accountsDTOs, nil
}
func (s *accountsService) UpdateByID(accountDTO *domain.AccountDTO) (*domain.AccountDTO, error) {
	account := accountDTO.ToAccount()

	account, err := s.repo.Update(account)
	if err != nil {
		return nil, err
	}

	return account.ToDTO(), nil
}
func (s *accountsService) Delete(id int, userID int) error {
	return s.repo.Delete(id, userID)
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
