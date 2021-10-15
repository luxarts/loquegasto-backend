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
