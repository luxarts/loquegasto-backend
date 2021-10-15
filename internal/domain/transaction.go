package domain

import (
	"time"
)

type Transaction struct {
	ID          string
	UserID      int
	MsgID       int
	Amount      float64
	Description string
	AccountID   *int
	CreatedAt   *time.Time
}

func (txn *Transaction) ToDTO() *TransactionDTO {
	dto := TransactionDTO{
		ID:          txn.ID,
		MsgID:       txn.MsgID,
		UserID:      txn.UserID,
		Amount:      txn.Amount,
		Description: txn.Description,
		AccountID:   txn.AccountID,
		CreatedAt:   txn.CreatedAt,
	}

	return &dto
}

type TransactionDTO struct {
	ID          string     `json:"id,omitempty"`
	MsgID       int        `json:"msg_id,omitempty"`
	UserID      int        `json:"user_id,omitempty"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	AccountID   *int       `json:"account_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
}

func (dto *TransactionDTO) IsValid() bool {
	return dto.Description != "" &&
		dto.Amount != 0 &&
		dto.MsgID != 0 &&
		dto.CreatedAt != nil
}

func (dto *TransactionDTO) IsValidForUpdate() bool {
	return dto.Description != "" &&
		dto.Amount != 0
}

func (dto *TransactionDTO) ToTransaction() *Transaction {
	txn := Transaction{
		ID:          dto.ID,
		MsgID:       dto.MsgID,
		Amount:      dto.Amount,
		UserID:      dto.UserID,
		Description: dto.Description,
		AccountID:   dto.AccountID,
		CreatedAt:   dto.CreatedAt,
	}

	return &txn
}
