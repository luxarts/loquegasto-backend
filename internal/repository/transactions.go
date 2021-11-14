package repository

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"
	"strings"

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
	UpdateByMsgID(transaction *domain.Transaction) (*domain.Transaction, error)
	GetAllByUserID(userID int, filters *domain.TransactionFilters) (*[]domain.Transaction, error)
	GetByMsgID(msgID int, userID int) (*domain.Transaction, error)
}

type transactionsRepository struct {
	db         *sqlx.DB
	sqlBuilder *transactionsSQL
}

func NewTransactionsRepository(db *sqlx.DB) TransactionsRepository {
	return &transactionsRepository{
		db:         db,
		sqlBuilder: &transactionsSQL{},
	}
}
func (r *transactionsRepository) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	transaction.ID = uuid.NewString()

	query, args, err := r.sqlBuilder.CreateSQL(transaction)
	if err != nil {
		return nil, jsend.NewError("failed CreateSQL", err, http.StatusInternalServerError)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}

	return transaction, nil
}
func (r *transactionsRepository) UpdateByMsgID(transaction *domain.Transaction) (*domain.Transaction, error) {
	query, args, err := r.sqlBuilder.UpdateByMsgIDSQL(transaction)
	if err != nil {
		return nil, jsend.NewError("failed UpdateByMsgIDSQL", err, http.StatusInternalServerError)
	}
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}
	return transaction, nil
}
func (r *transactionsRepository) GetAllByUserID(userID int, filters *domain.TransactionFilters) (*[]domain.Transaction, error) {
	query, args, err := r.sqlBuilder.GetAllByUserIDSQL(userID, filters)

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Queryx", err, http.StatusInternalServerError)
	}

	var results []domain.Transaction
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, jsend.NewError("failed Err: error in row", err, http.StatusInternalServerError)
		}
		var t domain.Transaction
		if err := rows.StructScan(&t); err != nil {
			return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
		}
		results = append(results, t)
	}

	return &results, nil
}
func (r *transactionsRepository) GetByMsgID(msgID int, userID int) (*domain.Transaction, error) {
	query, args, err := r.sqlBuilder.GetByMsgIDSQL(msgID, userID)
	if err != nil {
		return nil, jsend.NewError("failed GetByMsgIDSQL", err, http.StatusInternalServerError)
	}
	var transaction domain.Transaction
	err = r.db.QueryRowx(query, args...).StructScan(&transaction)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("transaction not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
	}

	return &transaction, nil
}

// SQL builders
type transactionsSQL struct{}

func (tsql *transactionsSQL) CreateSQL(transaction *domain.Transaction) (string, []interface{}, error) {
	return sq.Insert(tableTransactions).
		Columns(dbstruct.GetColumns(transaction)...).
		Values(dbstruct.GetValues(transaction)...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (tsql *transactionsSQL) UpdateByMsgIDSQL(transaction *domain.Transaction) (string, []interface{}, error) {
	return sq.Update(tableTransactions).
		Set("amount", transaction.Amount).
		Set("description", transaction.Description).
		Set("wallet_id", transaction.WalletID).
		Set("category_id", transaction.CategoryID).
		Where(sq.And{
			sq.Eq{"msg_id": transaction.MsgID},
			sq.Eq{"user_id": transaction.UserID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (tsql *transactionsSQL) GetAllByUserIDSQL(userID int, filters *domain.TransactionFilters) (string, []interface{}, error) {
	q := sq.Select("*").
		From(tableTransactions)

	if len(*filters) > 0 {
		and := sq.And{sq.Eq{"user_id": userID}}
		for k, v := range *filters {
			and = append(and, sq.Eq{k: v})
		}
		q = q.Where(and)
	} else {
		q = q.Where(sq.Eq{"user_id": userID})
	}

	return q.PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (tsql *transactionsSQL) GetByMsgIDSQL(msgID int, userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableTransactions).
		Where(sq.And{
			sq.Eq{"msg_id": msgID},
			sq.Eq{"user_id": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
