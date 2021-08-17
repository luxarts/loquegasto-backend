package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty"`
	MsgID       int                 `bson:"msg_id"`
	Amount      int64               `bson:"amount"`
	Description string              `bson:"description"`
	Source      string              `bson:"source,omitempty"`
	CreatedAt   time.Time           `bson:"created_at"`
}

func (txn *Transaction) ToDTO() *TransactionDTO {
	return &TransactionDTO{
		ID:          txn.ID.Hex(),
		MsgID:       txn.MsgID,
		Amount:      txn.Amount,
		Description: txn.Description,
		Source:      txn.Source,
		CreatedAt:   txn.CreatedAt,
	}
}

type TransactionDTO struct {
	ID          string    `json:"id,omitempty"`
	MsgID       int       `json:"msg_id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	Source      string    `json:"source,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

func (dto *TransactionDTO) IsValid() bool {
	return dto.Description != "" &&
		dto.Amount != 0 &&
		dto.MsgID != 0
}
func (dto *TransactionDTO) ToTransaction() *Transaction {
	txn := Transaction{
		MsgID:       dto.MsgID,
		Amount:      dto.Amount,
		Description: dto.Description,
		Source:      dto.Source,
	}

	objectID, err := primitive.ObjectIDFromHex(dto.ID)
	if err == nil {
		txn.ID = &objectID
	}

	return &txn
}
