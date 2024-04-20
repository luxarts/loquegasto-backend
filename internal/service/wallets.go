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
	GetByName(userID string, name string) (*domain.WalletCreateResponse, error)
	GetByID(userID string, id string) (*domain.WalletCreateResponse, error)
	GetAll(userID string) (*[]domain.WalletCreateResponse, error)
	UpdateByID(walletDTO *domain.WalletCreateRequest) (*domain.WalletCreateResponse, error)
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
		return nil, jsend.NewError("wallet name already exists", nil, http.StatusBadRequest)
	}

	wallet := req.ToWallet()
	wallet.ID = uuid.NewString()
	wallet.SanitizedName = sanitizedName
	wallet.UserID = userID
	wallet.CreatedAt = time.Now()

	wallet, err = s.repo.Create(wallet)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletCreateResponse(), nil
}
func (s *walletsService) GetByName(userID string, name string) (*domain.WalletCreateResponse, error) {
	wallet, err := s.repo.GetBySanitizedName(sanitizer.Sanitize(name), userID)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletCreateResponse(), nil
}
func (s *walletsService) GetByID(userID string, id string) (*domain.WalletCreateResponse, error) {
	wallet, err := s.repo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletCreateResponse(), nil
}
func (s *walletsService) GetAll(userID string) (*[]domain.WalletCreateResponse, error) {
	var err error
	var wallets *[]domain.Wallet

	wallets, err = s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var walletDTOs = make([]domain.WalletCreateResponse, len(*wallets))
	for _, wallet := range *wallets {
		walletDTOs = append(walletDTOs, *wallet.ToWalletCreateResponse())
	}

	return &walletDTOs, nil
}
func (s *walletsService) UpdateByID(walletDTO *domain.WalletCreateRequest) (*domain.WalletCreateResponse, error) {
	wallet := walletDTO.ToWallet()

	wallet.SanitizedName = sanitizer.Sanitize(wallet.Name)

	wallet, err := s.repo.UpdateByID(wallet)
	if err != nil {
		return nil, err
	}

	return wallet.ToWalletCreateResponse(), nil
}
func (s *walletsService) DeleteByID(id string, userID string) error {
	return s.repo.DeleteByID(id, userID)
}
