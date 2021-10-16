package repository

import (
	"loquegasto-backend/internal/domain"
	"net/http"

	"github.com/luxarts/jsend-go"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	tableWallets = "backend.wallets"
)

type WalletRepository interface {
	Create(account *domain.Wallet) (*domain.Wallet, error)
	GetAllByUserID(userID int) (*[]domain.Wallet, error)
	GetByID(id int) (*domain.Wallet, error)
	Update(account *domain.Wallet) (*domain.Wallet, error)
	Delete(id int, userID int) error
}
type walletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) WalletRepository {
	return &walletRepository{
		db: db,
	}
}
func (r *walletRepository) Create(wallet *domain.Wallet) (*domain.Wallet, error) {
	query := sq.Insert(tableWallets).Columns("user_id", "name", "balance", "created_at").
		Values(
			wallet.UserID,
			wallet.Name,
			wallet.Balance,
			wallet.CreatedAt).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&wallet.ID)
	if err != nil {
		return nil, jsend.NewError("failed QueryRow", err, http.StatusInternalServerError)
	}

	return wallet, nil
}
func (r *walletRepository) GetAllByUserID(userID int) (*[]domain.Wallet, error) {
	query, args, err := sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"user_id": userID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Query", err, http.StatusInternalServerError)
	}

	var results []domain.Wallet
	for rows.Next() {
		var wallet domain.Wallet
		if err := rows.StructScan(&wallet); err != nil {
			return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
		}
		results = append(results, wallet)
	}
	if err := rows.Err(); err != nil {
		return nil, jsend.NewError("failed Err", err, http.StatusInternalServerError)
	}

	return &results, nil
}
func (r *walletRepository) GetByID(id int) (*domain.Wallet, error) {
	query, args, err := sq.Select("*").
		From(tableWallets).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	var wallet domain.Wallet
	err = r.db.QueryRowx(query, args...).StructScan(&wallet)
	if err != nil {
		return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
	}

	return &wallet, nil
}
func (r *walletRepository) Update(wallet *domain.Wallet) (*domain.Wallet, error) {
	query, args, err := sq.Update(tableWallets).
		Set("id", wallet.ID).
		Set("user_id", wallet.UserID).
		Set("name", wallet.Name).
		Set("balance", wallet.Balance).
		Set("created_at", wallet.CreatedAt).
		Where(sq.Eq{"id": wallet.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	err = r.db.QueryRowx(query, args...).Err()
	if err != nil {
		return nil, jsend.NewError("failed QueryRowx", err, http.StatusInternalServerError)
	}
	return wallet, nil
}
func (r *walletRepository) Delete(id int, userID int) error {
	query, args, err := sq.Delete(tableWallets).
		Where(sq.Eq{"id": id, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	err = r.db.QueryRowx(query, args...).Err()
	if err != nil {
		return jsend.NewError("failed QueryRowx", err, http.StatusInternalServerError)
	}
	return nil
}
