package repository

import (
	"errors"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"
	"strings"

	"github.com/lib/pq"

	"github.com/luxarts/jsend-go"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	tableWallets = "core.wallets"
)

type WalletRepository interface {
	Create(wallet *domain.Wallet) (*domain.Wallet, error)
	GetAllByUserID(userID string) (*[]domain.Wallet, error)
	GetBySanitizedName(name string, userID string) (*domain.Wallet, error)
	GetByID(id string, userID string) (*domain.Wallet, error)
	UpdateByID(wallet *domain.Wallet, id string, userID string) (*domain.Wallet, error)
	DeleteByID(id string, userID string) error
}
type walletRepository struct {
	db         *sqlx.DB
	sqlBuilder *walletsSQL
}

func NewWalletRepository(db *sqlx.DB) WalletRepository {
	return &walletRepository{
		db:         db,
		sqlBuilder: &walletsSQL{},
	}
}
func (r *walletRepository) Create(wallet *domain.Wallet) (*domain.Wallet, error) {
	query, args, err := r.sqlBuilder.CreateSQL(wallet)
	if err != nil {
		return nil, jsend.NewError("failed walletRepository.Create.CreateSQL", err, http.StatusInternalServerError)
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			if pgerr.Code == defines.PGCodeDuplicateKey {
				return nil, jsend.NewError("wallet ID already exists", nil, http.StatusConflict)
			}
		}
		return nil, jsend.NewError("failed walletRepository.Create.Exec", err, http.StatusInternalServerError)
	}

	return wallet, nil
}
func (r *walletRepository) GetAllByUserID(userID string) (*[]domain.Wallet, error) {
	query, args, err := r.sqlBuilder.GetAllByUserIDSQL(userID)

	if err != nil {
		return nil, jsend.NewError("failed GetAllByUserIDSQL", err, http.StatusInternalServerError)
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Queryx", err, http.StatusInternalServerError)
	}

	var results []domain.Wallet
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, jsend.NewError("failed Err: error in row", err, http.StatusInternalServerError)
		}
		var wallet domain.Wallet
		if err := rows.StructScan(&wallet); err != nil {
			return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
		}
		results = append(results, wallet)
	}

	return &results, nil
}
func (r *walletRepository) GetBySanitizedName(name string, userID string) (*domain.Wallet, error) {
	query, args, err := r.sqlBuilder.GetBySanitizedNameSQL(name, userID)

	if err != nil {
		return nil, jsend.NewError("failed walletRepository.GetBySanitizedName.GetBySanitizedNameSQL", err, http.StatusInternalServerError)
	}

	var wallet domain.Wallet
	err = r.db.QueryRowx(query, args...).StructScan(&wallet)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("wallet not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed walletRepository.GetBySanitizedName.StructScan", err, http.StatusInternalServerError)
	}

	return &wallet, nil
}
func (r *walletRepository) GetByID(id string, userID string) (*domain.Wallet, error) {
	query, args, err := r.sqlBuilder.GetByIDSQL(id, userID)

	if err != nil {
		return nil, jsend.NewError("failed GetByIDSQL", err, http.StatusInternalServerError)
	}

	var wallet domain.Wallet
	err = r.db.QueryRowx(query, args...).StructScan(&wallet)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("wallet not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
	}

	return &wallet, nil
}
func (r *walletRepository) UpdateByID(wallet *domain.Wallet, id string, userID string) (*domain.Wallet, error) {
	query, args, err := r.sqlBuilder.UpdateByIDSQL(wallet, id, userID)

	if err != nil {
		return nil, jsend.NewError("failed walletRepository.UpdateByID.UpdateSQL", err, http.StatusInternalServerError)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed walletRepository.UpdateByID.Exec", err, http.StatusInternalServerError)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return nil, jsend.NewError("wallet not found", nil, http.StatusNotFound)
	}

	wallet.ID = id
	return wallet, nil
}
func (r *walletRepository) DeleteByID(id string, userID string) error {
	query, args, err := r.sqlBuilder.DeleteByIDSQL(id, userID)

	if err != nil {
		return jsend.NewError("failed DeleteByIDSQL", err, http.StatusInternalServerError)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return jsend.NewError("wallet not found", nil, http.StatusNotFound)
	}

	return nil
}

// SQL builders
type walletsSQL struct{}

func (wsql *walletsSQL) CreateSQL(wallet *domain.Wallet) (string, []interface{}, error) {
	return sq.Insert(tableWallets).
		Columns(dbstruct.GetColumns(wallet)...).
		Values(dbstruct.GetValues(wallet)...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) GetAllByUserIDSQL(userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) GetByIDSQL(id string, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"id": id, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) GetBySanitizedNameSQL(name string, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"sanitized_name": name, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) UpdateByIDSQL(wallet *domain.Wallet, id string, userID string) (string, []interface{}, error) {
	return sq.Update(tableWallets).
		Set("name", wallet.Name).
		Set("sanitized_name", wallet.SanitizedName).
		Set("balance", wallet.Balance).
		Set("emoji", wallet.Emoji).
		Where(sq.Eq{"id": id, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) DeleteByIDSQL(id string, userID string) (string, []interface{}, error) {
	return sq.Update(tableWallets).
		Set("deleted", true).
		Where(sq.Eq{"id": id, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
