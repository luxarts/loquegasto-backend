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
	GetAllByUserID(userID int) (*[]domain.Transaction, error)
	UpdateByMsgID(msgID int, transaction *domain.Transaction) (*domain.Transaction, error)
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
func (r *transactionsRepository) GetAllByUserID(userID int) (*[]domain.Transaction, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCtx()

	cursor, err := r.collection.Find(ctx, primitive.M{"user_id": userID})
	if err != nil {
		return nil, jsend.NewError("find-error", err)
	}
	var results []domain.Transaction
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, jsend.NewError("all-error", err)
	}

	return &results, nil
}
func (r *transactionsRepository) UpdateByMsgID(msgID int, transaction *domain.Transaction) (*domain.Transaction, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCtx()

	update := bson.D{
		{"$set", bson.D{
			{"amount", transaction.Amount},
			{"description", transaction.Description},
			{"source", transaction.Source},
		},
		},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"msg_id": msgID}, update)
	if err != nil {
		return nil, jsend.NewError("replaceone-error", err)
	}

	if res.MatchedCount != 1 {
		return nil, jsend.NewFail("not found")
	}

	return transaction, nil
}
