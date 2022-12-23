package repository

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"
	"strings"

	"github.com/luxarts/jsend-go"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	tableWallets = "core.wallets"
)

type WalletRepository interface {
	Create(wallet *domain.Wallet) (*domain.Wallet, error)
	GetAllByUserID(userID int) (*[]domain.Wallet, error)
	GetBySanitizedName(name string, userID int) (*domain.Wallet, error)
	GetByID(id int, userID int) (*domain.Wallet, error)
	UpdateByID(wallet *domain.Wallet) (*domain.Wallet, error)
	DeleteByID(id int, userID int) error
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
		return nil, jsend.NewError("failed CreateSQL", err, http.StatusInternalServerError)
	}

	err = r.db.QueryRowx(query, args...).Scan(&wallet.ID)
	if err != nil {
		return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
	}

	return wallet, nil
}
func (r *walletRepository) GetAllByUserID(userID int) (*[]domain.Wallet, error) {
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
func (r *walletRepository) GetBySanitizedName(name string, userID int) (*domain.Wallet, error) {
	query, args, err := r.sqlBuilder.GetBySanitizedNameSQL(name, userID)

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
func (r *walletRepository) GetByID(id int, userID int) (*domain.Wallet, error) {
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
func (r *walletRepository) UpdateByID(wallet *domain.Wallet) (*domain.Wallet, error) {
	query, args, err := r.sqlBuilder.UpdateByIDSQL(wallet)

	if err != nil {
		return nil, jsend.NewError("failed UpdateSQL", err, http.StatusInternalServerError)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return nil, jsend.NewError("wallet not found", nil, http.StatusNotFound)
	}

	return wallet, nil
}
func (r *walletRepository) DeleteByID(id int, userID int) error {
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
	wallet.ID = nil
	return sq.Insert(tableWallets).
		Columns(dbstruct.GetColumns(wallet)[1:]...).
		Values(dbstruct.GetValues(wallet)[1:]...).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) GetAllByUserIDSQL(userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) GetByIDSQL(id int, userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"id": id, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) GetBySanitizedNameSQL(name string, userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"sanitized_name": name, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) UpdateByIDSQL(wallet *domain.Wallet) (string, []interface{}, error) {
	return sq.Update(tableWallets).
		Set("name", wallet.Name).
		Set("balance", wallet.Balance).
		Where(sq.Eq{"id": wallet.ID, "user_id": wallet.UserID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (wsql *walletsSQL) DeleteByIDSQL(id int, userID int) (string, []interface{}, error) {
	return sq.Delete(tableWallets).
		Where(sq.Eq{"id": id, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
