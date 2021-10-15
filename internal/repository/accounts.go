package repository

import (
	"loquegasto-backend/internal/domain"
	"net/http"

	"github.com/luxarts/jsend-go"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	tableAccounts = "backend.accounts"
)

type AccountRepository interface {
	Create(account *domain.Account) (*domain.Account, error)
	GetAllByUserID(userID int) (*[]domain.Account, error)
}
type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}
func (r *accountRepository) Create(account *domain.Account) (*domain.Account, error) {
	query := sq.Insert(tableAccounts).Columns("user_id", "name", "balance", "updated_at").
		Values(
			account.UserID,
			account.Name,
			account.Balance,
			account.UpdatedAt).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&account.ID)
	if err != nil {
		return nil, jsend.NewError("failed QueryRow", err, http.StatusInternalServerError)
	}

	return account, nil
}
func (r *accountRepository) GetAllByUserID(userID int) (*[]domain.Account, error) {
	query := sq.Select("*").
		From(tableAccounts).
		Where(sq.Eq{"user_id": userID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.Query()
	if err != nil {
		return nil, jsend.NewError("failed Query", err, http.StatusInternalServerError)
	}

	var results []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Balance, &a.UpdatedAt); err != nil {
			return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
		}
		results = append(results, a)
	}
	if err := rows.Err(); err != nil {
		return nil, jsend.NewError("failed Err", err, http.StatusInternalServerError)
	}

	return &results, nil
}
