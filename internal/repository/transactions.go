package repository

import (
	"context"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"time"

	"github.com/luxarts/jsend-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionsRepository interface {
	Create(transaction *domain.Transaction) (*domain.Transaction, error)
}

type transactionsRepository struct {
	collection *mongo.Collection
}

func NewTransactionsRepository(client *mongo.Client) TransactionsRepository {
	return &transactionsRepository{
		collection: client.Database(defines.MongoDatabase).Collection(defines.MongoTransactionsCollection),
	}
}

func (r *transactionsRepository) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	transaction.CreatedAt = time.Now().UTC()
	transactionBson, err := bson.Marshal(transaction)
	if err != nil {
		return nil, jsend.NewError("marshal-error", err)
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCtx()

	resp, err := r.collection.InsertOne(ctx, transactionBson)
	if err != nil {
		return nil, jsend.NewError("insertone-error", err)
	}

	txnID := resp.InsertedID.(primitive.ObjectID)
	transaction.ID = &txnID

	return transaction, nil
}
