package repository

import (
	"fmt"
	"loquegasto-backend/internal/domain"
	"os"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TransactionsRepository interface {
	Create(transaction *domain.Transaction) (*domain.Transaction, error)
	GetAllByUserID(userID int) (*[]domain.Transaction, error)
	UpdateByMsgID(msgID int, transaction *domain.Transaction) (*domain.Transaction, error)
}

type transactionsRepositoryPostgreSQL struct {
	db *sqlx.DB
}

func NewTransactionsRepository() TransactionsRepository {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Sprintf("Fail to connect to database: %v", err))
	}

	return &transactionsRepositoryPostgreSQL{
		db: db,
	}
}
func (r *transactionsRepositoryPostgreSQL) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	id := uuid.NewString()

	query := sq.Insert("backend.transactions").Columns("uuid", "user_id", "msg_id", "amount", "description", "account_id", "created_at").
		Values(
			id,
			transaction.UserID,
			transaction.MsgID,
			transaction.Amount,
			transaction.Description,
			transaction.AccountID,
			transaction.CreatedAt).
		Suffix("RETURNING \"uuid\"").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed QueryRow: %w", err)
	}

	transaction.ID = id

	return transaction, nil
}
func (r *transactionsRepositoryPostgreSQL) GetAllByUserID(userID int) (*[]domain.Transaction, error) {
	return nil, nil
}
func (r *transactionsRepositoryPostgreSQL) UpdateByMsgID(msgID int, transaction *domain.Transaction) (*domain.Transaction, error) {
	return nil, nil
}
