package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
	"loquegasto-backend/internal/utils/sanitizer"
	"net/http"

	"github.com/luxarts/jsend-go"
)

type WalletsService interface {
	Create(accountDTO *domain.WalletDTO) (*domain.WalletDTO, error)
	GetByName(userID int, name string) (*domain.WalletDTO, error)
	GetByID(userID int, id int) (*domain.WalletDTO, error)
	GetAll(userID int, search string) (*[]domain.WalletDTO, error)
	UpdateByID(accountDTO *domain.WalletDTO) (*domain.WalletDTO, error)
	DeleteByID(id int, userID int) error
}
type walletsService struct {
	repo repository.WalletRepository
}

func NewWalletsService(repo repository.WalletRepository) WalletsService {
	return &walletsService{
		repo: repo,
	}
}
func (s *walletsService) Create(walletDTO *domain.WalletDTO) (*domain.WalletDTO, error) {
	wallet := walletDTO.ToWallet()

	wallet.SanitizedName = sanitizer.Sanitize(walletDTO.Name)

	wallet, err := s.repo.Create(wallet)
	if err != nil {
		return nil, err
	}

	return wallet.ToDTO(), nil
}
func (s *walletsService) GetByName(userID int, name string) (*domain.WalletDTO, error) {
	wallet, err := s.repo.GetBySanitizedName(userID, sanitizer.Sanitize(name))
	if err != nil {
		return nil, err
	}

	return wallet.ToDTO(), nil
}
func (s *walletsService) GetByID(userID int, id int) (*domain.WalletDTO, error) {
	wallet, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != userID {
		return nil, jsend.NewError("forbidden", nil, http.StatusForbidden)
	}

	return wallet.ToDTO(), nil
}
func (s *walletsService) GetAll(userID int, search string) (*[]domain.WalletDTO, error) {
	var err error
	var wallets *[]domain.Wallet

	if search != "" {
		var w *domain.Wallet
		w, err = s.repo.GetBySanitizedName(userID, sanitizer.Sanitize(search))
		if err != nil {
			return nil, err
		}
		wallets = &[]domain.Wallet{*w}
	} else {
		wallets, err = s.repo.GetAllByUserID(userID)
		if err != nil {
			return nil, err
		}
	}

	var walletDTOs = make([]domain.WalletDTO, 0)
	for _, wallet := range *wallets {
		walletDTOs = append(walletDTOs, *wallet.ToDTO())
	}

	return &walletDTOs, nil
}
func (s *walletsService) UpdateByID(walletDTO *domain.WalletDTO) (*domain.WalletDTO, error) {
	wallet := walletDTO.ToWallet()

	wallet, err := s.repo.UpdateByID(wallet)
	if err != nil {
		return nil, err
	}

	return wallet.ToDTO(), nil
}
func (s *walletsService) DeleteByID(id int, userID int) error {
	return s.repo.DeleteByID(id, userID)
}
