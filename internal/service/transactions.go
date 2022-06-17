package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type TransactionsService interface {
	Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
	UpdateByMsgID(userID int, transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
	GetAllByUserID(userID int, filters *domain.TransactionFilters) (*[]domain.TransactionDTO, error)
}

type transactionsService struct {
	txnRepo    repository.TransactionsRepository
	walletRepo repository.WalletRepository
}

func NewTransactionsService(txnRepo repository.TransactionsRepository, walletRepo repository.WalletRepository) TransactionsService {
	return &transactionsService{
		txnRepo:    txnRepo,
		walletRepo: walletRepo,
	}
}

func (s *transactionsService) Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	transaction := transactionDTO.ToTransaction()

	transaction, err := s.txnRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Update balance
	wallet, err := s.walletRepo.GetByID(transaction.WalletID, transaction.UserID)
	if err != nil {
		return nil, err
	}
	wallet.Balance += transaction.Amount
	wallet, err = s.walletRepo.UpdateByID(wallet)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) UpdateByMsgID(userID int, transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	currentTransaction, err := s.txnRepo.GetByMsgID(transactionDTO.MsgID, userID)
	if err != nil {
		return nil, err
	}

	transactionDTO.UserID = userID
	transaction := transactionDTO.ToTransaction()

	transaction, err = s.txnRepo.UpdateByMsgID(transaction)
	if err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.GetByID(transaction.WalletID, transaction.UserID)
	if err != nil {
		return nil, err
	}

	// Rollback old transaction's amount
	wallet.Balance -= currentTransaction.Amount

	// Update balance with new transaction's amount
	wallet.Balance += transaction.Amount

	wallet, err = s.walletRepo.UpdateByID(wallet)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) GetAllByUserID(userID int, filters *domain.TransactionFilters) (*[]domain.TransactionDTO, error) {
	res, err := s.txnRepo.GetAllByUserID(userID, filters)
	if err != nil {
		return nil, err
	}

	var dtos []domain.TransactionDTO
	for _, r := range *res {
		dtos = append(dtos, *r.ToDTO())
	}

	return &dtos, nil
}
