package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/luxarts/jsend-go"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/utils/sanitizer"
	"net/http"
	"time"
)

type WalletsService interface {
	Create(req *domain.WalletCreateRequest, userID string) (*domain.WalletCreateResponse, error)
	GetByID(id string, userID string) (*domain.WalletGetResponse, error)
	GetAll(userID string) (*[]domain.WalletGetResponse, error)
	UpdateByID(req *domain.WalletUpdateRequest, id string, userID string) (*domain.WalletUpdateResponse, error)
	DeleteByID(id string, userID string) error
}
type walletsService struct {
	repo repository.WalletRepository
}

func NewWalletsService(repo repository.WalletRepository) WalletsService {
	return &walletsService{
		repo: repo,
	}
}
func (s *walletsService) Create(req *domain.WalletCreateRequest, userID string) (*domain.WalletCreateResponse, error) {
	sanitizedName := sanitizer.Sanitize(req.Name)

	// Check if the name already exists for the given user
	w, err := s.repo.GetBySanitizedName(sanitizedName, userID)
	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && err != nil && *jsendErr.Code != http.StatusNotFound {
		return nil, err
	}
	if w != nil {
		return nil, jsend.NewError("wallet name already used", nil, http.StatusBadRequest)
	}

	wallet := req.ToWallet()
	wallet.ID = uuid.NewString()
	wallet.UserID = userID
	wallet.CreatedAt = time.Now()

	wallet, err = s.repo.Create(wallet)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletCreateResponse(), nil
}
func (s *walletsService) GetByName(name string, userID string) (*domain.WalletGetResponse, error) {
	wallet, err := s.repo.GetBySanitizedName(sanitizer.Sanitize(name), userID)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletGetResponse(), nil
}
func (s *walletsService) GetByID(id string, userID string) (*domain.WalletGetResponse, error) {
	wallet, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletGetResponse(), nil
}
func (s *walletsService) GetAll(userID string) (*[]domain.WalletGetResponse, error) {
	wallets, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response = make([]domain.WalletGetResponse, len(*wallets))
	for i, wallet := range *wallets {
		response[i] = *wallet.ToWalletGetResponse()
	}

	return &response, nil
}
func (s *walletsService) UpdateByID(req *domain.WalletUpdateRequest, id string, userID string) (*domain.WalletUpdateResponse, error) {
	sanitizedName := sanitizer.Sanitize(req.Name)

	// Check if the name already exists for the given user
	w, err := s.repo.GetBySanitizedName(sanitizedName, userID)
	var jsendErr *jsend.Body
	if errors.As(err, &jsendErr) && err != nil && *jsendErr.Code != http.StatusNotFound {
		return nil, err
	}
	if w != nil {
		return nil, jsend.NewError("wallet name already used", nil, http.StatusBadRequest)
	}

	wallet := req.ToWallet()

	wallet, err = s.repo.UpdateByID(wallet, id, userID)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletUpdateResponse(), nil
}
func (s *walletsService) DeleteByID(id string, userID string) error {
	return s.repo.DeleteByID(id, userID)
}
