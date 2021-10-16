package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type TransactionsService interface {
	Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
	UpdateByMsgID(userID int, msgID int, transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
	GetAllByUserID(userID int) (*[]domain.TransactionDTO, error)
}

type transactionsService struct {
	txnRepo repository.TransactionsRepository
	accRepo repository.AccountRepository
}

func NewTransactionsService(txnRepo repository.TransactionsRepository, accRepo repository.AccountRepository) TransactionsService {
	return &transactionsService{
		txnRepo: txnRepo,
		accRepo: accRepo,
	}
}

func (s *transactionsService) Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	transaction := transactionDTO.ToTransaction()

	transaction, err := s.txnRepo.Create(transaction)
	if err != nil {
		return nil, err
	}

	// Update balance
	account, err := s.accRepo.GetByID(transaction.AccountID)
	if err != nil {
		return nil, err
	}
	account.Balance -= int64(transaction.Amount * 100)
	account, err = s.accRepo.Update(account)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) UpdateByMsgID(userID int, msgID int, transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	transactionDTO.UserID = userID

	transaction := transactionDTO.ToTransaction()

	transaction, err := s.txnRepo.UpdateByMsgID(msgID, transaction)
	if err != nil {
		return nil, err
	}

	// Update balance
	account, err := s.accRepo.GetByID(transaction.AccountID)
	if err != nil {
		return nil, err
	}
	account.Balance -= int64(transaction.Amount * 100)
	account, err = s.accRepo.Update(account)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) GetAllByUserID(userID int) (*[]domain.TransactionDTO, error) {
	res, err := s.txnRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var dtos []domain.TransactionDTO
	for _, r := range *res {
		dtos = append(dtos, *r.ToDTO())
	}

	return &dtos, nil
}
