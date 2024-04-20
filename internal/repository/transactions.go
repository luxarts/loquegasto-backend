package repository

import (
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	tableTransactions = "core.transactions"
)

type TransactionsRepository interface {
	Create(transaction *domain.Transaction) (*domain.Transaction, error)
	UpdateByMsgID(transaction *domain.Transaction) (*domain.Transaction, error)
	GetAll(userID int64, filters *domain.TransactionFilters) (*[]domain.Transaction, error)
	GetByMsgID(msgID int64, userID string) (*domain.Transaction, error)
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
func (r *transactionsRepository) Create(t *domain.Transaction) (*domain.Transaction, error) {
	t.ID = uuid.NewString()

	query, args, err := r.sqlBuilder.CreateSQL(t)
	if err != nil {
		return nil, jsend.NewError("failed CreateSQL", err, http.StatusInternalServerError)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}

	return t, nil
}
func (r *transactionsRepository) UpdateByMsgID(t *domain.Transaction) (*domain.Transaction, error) {
	query, args, err := r.sqlBuilder.UpdateByMsgIDSQL(t)
	if err != nil {
		return nil, jsend.NewError("failed UpdateByMsgIDSQL", err, http.StatusInternalServerError)
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, jsend.NewError("failed RowsAffected", err, http.StatusInternalServerError)
	}
	if affected == 0 {
		return nil, jsend.NewError("transaction not found", nil, http.StatusNotFound)
	}

	return t, nil
}
func (r *transactionsRepository) GetAll(userID int64, filters *domain.TransactionFilters) (*[]domain.Transaction, error) {
	query, args, err := r.sqlBuilder.GetAllSQL(userID, filters)

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
func (r *transactionsRepository) GetByMsgID(msgID int64, userID string) (*domain.Transaction, error) {
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

func (tsql *transactionsSQL) CreateSQL(t *domain.Transaction) (string, []interface{}, error) {
	return sq.Insert(tableTransactions).
		Columns(dbstruct.GetColumns(t)...).
		Values(dbstruct.GetValues(t)...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (tsql *transactionsSQL) UpdateByMsgIDSQL(t *domain.Transaction) (string, []interface{}, error) {
	builder := sq.Update(tableTransactions)

	builder = dbstruct.SetValues(builder, t)

	return builder.
		Where(sq.And{
			sq.Eq{"msg_id": t.MsgID},
			sq.Eq{"user_id": t.UserID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (tsql *transactionsSQL) GetAllSQL(userID int64, filters *domain.TransactionFilters) (string, []interface{}, error) {
	q := sq.Select("*").
		From(tableTransactions)

	if len(*filters) > 0 {
		and := sq.And{sq.Eq{"user_id": userID}}
		for k, v := range *filters {
			if k == defines.QueryFrom || k == defines.QueryTo {
				tsInt, err := strconv.ParseInt(v, 10, 64)
				if err == nil {
					ts := time.Unix(tsInt, 0).Format(time.RFC3339)
					if k == defines.QueryFrom {
						and = append(and, sq.GtOrEq{"created_at": ts})
					} else {
						and = append(and, sq.LtOrEq{"created_at": ts})
					}
				}
				continue
			}

			and = append(and, sq.Eq{k: v})
		}
		q = q.Where(and)
	} else {
		q = q.Where(sq.Eq{"user_id": userID})
	}
	q = q.OrderBy("created_at DESC")
	return q.PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (tsql *transactionsSQL) GetByMsgIDSQL(msgID int64, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableTransactions).
		Where(sq.And{
			sq.Eq{"msg_id": msgID},
			sq.Eq{"user_id": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
