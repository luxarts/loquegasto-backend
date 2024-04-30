package domain

import (
	"loquegasto-backend/internal/defines"
	"time"
)

type Transaction struct {
	ID          string    `db:"uuid"`
	UserID      string    `db:"user_id"`
	Amount      int64     `db:"amount"`
	Description string    `db:"description"`
	WalletID    string    `db:"wallet_id"`
	CategoryID  string    `db:"category_id"`
	MsgID       *int64    `db:"msg_id"`
	CreatedAt   time.Time `db:"created_at"`
}

func (txn *Transaction) ToTransactionCreateResponse() *TransactionCreateResponse {
	return &TransactionCreateResponse{
		ID:          txn.ID,
		Amount:      float64(txn.Amount) / 100.0,
		Description: txn.Description,
		WalletID:    txn.WalletID,
		CreatedAt:   txn.CreatedAt,
		CategoryID:  txn.CategoryID,
	}
}

type TransactionCreateRequest struct {
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	WalletID    string     `json:"wallet_id"`
	CategoryID  string     `json:"category_id"`
	CreatedAt   *time.Time `json:"created_at"`
	MsgID       *int64     `json:"msg_id"`
}

type TransactionCreateResponse struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	WalletID    string    `json:"wallet_id"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (req *TransactionCreateRequest) IsValid() bool {
	return req.Description != "" &&
		req.Amount != 0 &&
		req.WalletID != "" &&
		req.CategoryID != "" &&
		req.CreatedAt != nil
}

func (req *TransactionCreateRequest) IsValidForUpdate() bool {
	return req.Description != "" &&
		req.Amount != 0
}

func (req *TransactionCreateRequest) ToTransaction() *Transaction {
	return &Transaction{
		MsgID:       req.MsgID,
		Amount:      int64(req.Amount * 100),
		Description: req.Description,
		WalletID:    req.WalletID,
		CreatedAt:   *req.CreatedAt,
		CategoryID:  req.CategoryID,
	}
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
