package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type TransactionsService interface {
	Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
	UpdateByMsgID(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
	GetAll(userID int, filters *domain.TransactionFilters) (*[]domain.TransactionDTO, error)
}

type transactionsService struct {
	txnRepo    repository.TransactionsRepository
	walletRepo repository.WalletRepository
	catRepo    repository.CategoriesRepository
}

func NewTransactionsService(txnRepo repository.TransactionsRepository, walletRepo repository.WalletRepository, catRepo repository.CategoriesRepository) TransactionsService {
	return &transactionsService{
		txnRepo:    txnRepo,
		walletRepo: walletRepo,
		catRepo:    catRepo,
	}
}

func (s *transactionsService) Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	// Check if wallet exists
	wallet, err := s.walletRepo.GetByID(transactionDTO.WalletID, transactionDTO.UserID)
	if err != nil {
		return nil, err
	}

	// Check if category exists
	if transactionDTO.CategoryID != nil {
		_, err = s.catRepo.GetByID(*transactionDTO.CategoryID, transactionDTO.UserID)
		if err != nil {
			return nil, err
		}
	}

	transaction := transactionDTO.ToTransaction()

	transaction, err = s.txnRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Update balance
	wallet.Balance += transaction.Amount
	wallet, err = s.walletRepo.UpdateByID(wallet)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) UpdateByMsgID(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	// Check if wallet exists
	walletDest, err := s.walletRepo.GetByID(transactionDTO.WalletID, transactionDTO.UserID)
	if err != nil {
		return nil, err
	}

	// Check if category exists
	if transactionDTO.CategoryID != nil {
		_, err = s.catRepo.GetByID(*transactionDTO.CategoryID, transactionDTO.UserID)
		if err != nil {
			return nil, err
		}
	}

	currentTransaction, err := s.txnRepo.GetByMsgID(transactionDTO.MsgID, transactionDTO.UserID)
	if err != nil {
		return nil, err
	}

	walletSrc, err := s.walletRepo.GetByID(currentTransaction.WalletID, transactionDTO.UserID)
	if err != nil {
		return nil, err
	}

	transaction := transactionDTO.ToTransaction()

	transaction, err = s.txnRepo.UpdateByMsgID(transaction)
	if err != nil {
		return nil, err
	}

	// Rollback old transaction's amount
	walletSrc.Balance -= currentTransaction.Amount

	// Update balance with new transaction's amount
	walletDest.Balance += transaction.Amount

	_, err = s.walletRepo.UpdateByID(walletSrc)
	if err != nil {
		return nil, err
	}

	walletDest, err = s.walletRepo.UpdateByID(walletDest)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) GetAll(userID int, filters *domain.TransactionFilters) (*[]domain.TransactionDTO, error) {
	res, err := s.txnRepo.GetAll(userID, filters)
	if err != nil {
		return nil, err
	}

	dtos := make([]domain.TransactionDTO, 0)
	for _, r := range *res {
		dtos = append(dtos, *r.ToDTO())
	}

	return &dtos, nil
}
