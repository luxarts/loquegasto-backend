package service

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/repository"
)

type TransactionsService interface {
	Create(transactionDTO *domain.TransactionDTO) (*domain.TransactionDTO, error)
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
