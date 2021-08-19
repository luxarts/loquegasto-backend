package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty"`
	MsgID       int                 `bson:"msg_id,omitempty"`
	UserID      int                 `bson:"user_id,omitempty"`
	Amount      float64             `bson:"amount,omitempty"`
	Description string              `bson:"description,omitempty"`
	Source      string              `bson:"source,omitempty"`
	CreatedAt   *time.Time          `bson:"created_at,omitempty"`
}

func (txn *Transaction) ToDTO() *TransactionDTO {
	dto := TransactionDTO{
		MsgID:       txn.MsgID,
		UserID:      txn.UserID,
		Amount:      txn.Amount,
		Description: txn.Description,
		Source:      txn.Source,
		CreatedAt:   txn.CreatedAt,
	}

	if txn.ID != nil {
		dto.ID = txn.ID.Hex()
	}

	return &dto
}

type TransactionDTO struct {
	ID          string     `json:"id,omitempty"`
	MsgID       int        `json:"msg_id,omitempty"`
	UserID      int        `json:"user_id,omitempty"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	Source      string     `json:"source,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
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
		MsgID:       dto.MsgID,
		Amount:      dto.Amount,
		UserID:      dto.UserID,
		Description: dto.Description,
		Source:      dto.Source,
		CreatedAt:   dto.CreatedAt,
	}

	objectID, err := primitive.ObjectIDFromHex(dto.ID)
	if err == nil {
		txn.ID = &objectID
	}

	return &txn
}

type TotalDTO struct {
	Total float64 `json:"total"`
}
