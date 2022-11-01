package domain

import (
	"loquegasto-backend/internal/defines"
	"time"
)

type Transaction struct {
	ID          string     `db:"uuid"`
	UserID      int        `db:"user_id"`
	MsgID       int        `db:"msg_id"`
	Amount      int64      `db:"amount"`
	Description string     `db:"description"`
	WalletID    int        `db:"wallet_id"`
	CreatedAt   *time.Time `db:"created_at"`
	CategoryID  *int       `db:"category_id"`
}

func (txn *Transaction) ToDTO() *TransactionDTO {
	dto := TransactionDTO{
		ID:          txn.ID,
		MsgID:       txn.MsgID,
		UserID:      txn.UserID,
		Amount:      float64(txn.Amount) / 100.0,
		Description: txn.Description,
		WalletID:    txn.WalletID,
		CreatedAt:   txn.CreatedAt,
		CategoryID:  txn.CategoryID,
	}

	return &dto
}

type TransactionDTO struct {
	ID          string     `json:"id,omitempty"`
	MsgID       int        `json:"msg_id,omitempty"`
	UserID      int        `json:"user_id,omitempty"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	WalletID    int        `json:"wallet_id"`
	CreatedAt   *time.Time `json:"created_at"`
	CategoryID  *int       `json:"category_id,omitempty"`
}

func (dto *TransactionDTO) IsValid() bool {
	return dto.Description != "" &&
		dto.Amount != 0 &&
		dto.MsgID != 0 &&
		dto.WalletID != 0 &&
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
		Amount:      int64(dto.Amount * 100),
		UserID:      dto.UserID,
		Description: dto.Description,
		WalletID:    dto.WalletID,
		CreatedAt:   dto.CreatedAt,
		CategoryID:  dto.CategoryID,
	}

	return &txn
}

type TransactionFilters map[string]string

func (tf *TransactionFilters) IsValid() bool {
	for key := range *tf {
		if key != defines.QueryWalletID &&
			key != defines.QueryCategoryID {
			return false
		}
	}
	return true
}
