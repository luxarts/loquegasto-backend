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
	repo repository.TransactionsRepository
}

func NewTransactionsService(repo repository.TransactionsRepository) TransactionsService {
	return &transactionsService{
		repo: repo,
	}
}

func (s *transactionsService) Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	transaction := transactionDTO.ToTransaction()

	transaction, err := s.repo.Create(transaction)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) UpdateByMsgID(userID int, msgID int, transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	transactionDTO.UserID = userID

	transaction := transactionDTO.ToTransaction()

	transaction, err := s.repo.UpdateByMsgID(msgID, transaction)
	if err != nil {
		return nil, err
	}

	response := transaction.ToDTO()

	return response, nil
}
func (s *transactionsService) GetAllByUserID(userID int) (*[]domain.TransactionDTO, error) {
	res, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var dtos []domain.TransactionDTO
	for _, r := range *res {
		dtos = append(dtos, *r.ToDTO())
	}

	return &dtos, nil
}
