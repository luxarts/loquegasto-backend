package repository

import (
	"loquegasto-backend/internal/domain"
	"net/http"

	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	tableTransactions = "backend.transactions"
)

type TransactionsRepository interface {
	Create(transaction *domain.Transaction) (*domain.Transaction, error)
	UpdateByMsgID(msgID int, transaction *domain.Transaction) (*domain.Transaction, error)
	GetAllByUserID(userID int) (*[]domain.Transaction, error)
}

type transactionsRepository struct {
	db *sqlx.DB
}

func NewTransactionsRepository(db *sqlx.DB) TransactionsRepository {
	return &transactionsRepository{
		db: db,
	}
}
func (r *transactionsRepository) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	id := uuid.NewString()

	query := sq.Insert(tableTransactions).Columns("uuid", "user_id", "msg_id", "amount", "description", "account_id", "created_at").
		Values(
			id,
			transaction.UserID,
			transaction.MsgID,
			transaction.Amount,
			transaction.Description,
			transaction.WalletID,
			transaction.CreatedAt).
		Suffix("RETURNING \"uuid\"").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&id)
	if err != nil {
		return nil, jsend.NewError("failed QueryRow", err, http.StatusInternalServerError)
	}

	transaction.ID = id

	return transaction, nil
}
func (r *transactionsRepository) UpdateByMsgID(msgID int, transaction *domain.Transaction) (*domain.Transaction, error) {
	var id string

	query := sq.Update(tableTransactions).
		Set("amount", transaction.Amount).
		Set("description", transaction.Description).
		Set("account_id", transaction.WalletID).
		Where(sq.Eq{"msg_id": msgID}).
		Suffix("RETURNING \"uuid\"").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&id)
	if err != nil {
		return nil, jsend.NewError("failed QueryRow", err, http.StatusInternalServerError)
	}
	return transaction, nil
}
func (r *transactionsRepository) GetAllByUserID(userID int) (*[]domain.Transaction, error) {
	query := sq.Select("*").
		From(tableTransactions).
		Where(sq.Eq{"user_id": userID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.Query()
	if err != nil {
		return nil, jsend.NewError("failed Query", err, http.StatusInternalServerError)
	}

	var results []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.MsgID, &t.Amount, &t.Description, &t.WalletID, &t.CreatedAt); err != nil {
			return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
		}
		results = append(results, t)
	}
	if err := rows.Err(); err != nil {
		return nil, jsend.NewError("failed Err", err, http.StatusInternalServerError)
	}

	return &results, nil
}
