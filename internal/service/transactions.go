package service

import (
	"github.com/google/uuid"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type TransactionsService interface {
	Create(req *domain.TransactionCreateRequest, userID string) (*domain.TransactionCreateResponse, error)
	UpdateByMsgID(req *domain.TransactionCreateRequest, userID string) (*domain.TransactionCreateResponse, error)
	GetAll(filters *domain.TransactionFilters, userID string) (*[]domain.TransactionCreateResponse, error)
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

func (s *transactionsService) Create(req *domain.TransactionCreateRequest, userID string) (*domain.TransactionCreateResponse, error) {
	// Check if wallet exists
	wallet, err := s.walletRepo.GetByID(req.WalletID, userID)
	if err != nil {
		return nil, err
	}

	// Check if category exists
	_, err = s.catRepo.GetByID(req.CategoryID, userID)
	if err != nil {
		return nil, err
	}

	transaction := req.ToTransaction()
	transaction.ID = uuid.NewString()
	transaction.UserID = userID

	transaction, err = s.txnRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Update balance
	wallet.Balance += transaction.Amount
	wallet, err = s.walletRepo.UpdateByID(wallet, wallet.ID, userID)
	if err != nil {
		return nil, err
	}

	response := transaction.ToTransactionCreateResponse()

	return response, nil
}
func (s *transactionsService) UpdateByMsgID(transactionDTO *domain.TransactionCreateRequest, userID string) (*domain.TransactionCreateResponse, error) {
	// Check if wallet exists
	walletDest, err := s.walletRepo.GetByID(transactionDTO.WalletID, userID)
	if err != nil {
		return nil, err
	}

	_, err = s.catRepo.GetByID(transactionDTO.CategoryID, userID)
	if err != nil {
		return nil, err
	}

	currentTransaction, err := s.txnRepo.GetByMsgID(*transactionDTO.MsgID, userID)
	if err != nil {
		return nil, err
	}

	walletSrc, err := s.walletRepo.GetByID(currentTransaction.WalletID, userID)
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

	_, err = s.walletRepo.UpdateByID(walletSrc, walletSrc.ID, userID)
	if err != nil {
		return nil, err
	}

	walletDest, err = s.walletRepo.UpdateByID(walletDest, walletDest.ID, userID)
	if err != nil {
		return nil, err
	}

	response := transaction.ToTransactionCreateResponse()

	return response, nil
}
func (s *transactionsService) GetAll(filters *domain.TransactionFilters, userID string) (*[]domain.TransactionCreateResponse, error) {
	res, err := s.txnRepo.GetAll(userID, filters)
	if err != nil {
		return nil, err
	}

	txns := make([]domain.TransactionCreateResponse, 0)
	for _, r := range *res {
		txns = append(txns, *r.ToTransactionCreateResponse())
	}

	return &txns, nil
}
