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
	GetByID(id int) (*domain.Account, error)
	Update(account *domain.Account) (*domain.Account, error)
	Delete(id int, userID int) error
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
	query := sq.Insert(tableAccounts).Columns("user_id", "name", "balance", "created_at").
		Values(
			account.UserID,
			account.Name,
			account.Balance,
			account.CreatedAt).
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
	query, args, err := sq.Select("*").
		From(tableAccounts).
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

	var results []domain.Account
	for rows.Next() {
		var a domain.Account
		if err := rows.StructScan(&a); err != nil {
			return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
		}
		results = append(results, a)
	}
	if err := rows.Err(); err != nil {
		return nil, jsend.NewError("failed Err", err, http.StatusInternalServerError)
	}

	return &results, nil
}
func (r *accountRepository) GetByID(id int) (*domain.Account, error) {
	query, args, err := sq.Select("*").
		From(tableAccounts).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	var account domain.Account
	err = r.db.QueryRowx(query, args...).StructScan(&account)
	if err != nil {
		return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
	}

	return &account, nil
}
func (r *accountRepository) Update(account *domain.Account) (*domain.Account, error) {
	query, args, err := sq.Update(tableAccounts).
		Set("id", account.ID).
		Set("user_id", account.UserID).
		Set("name", account.Name).
		Set("balance", account.Balance).
		Set("created_at", account.CreatedAt).
		Where(sq.Eq{"id": account.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	err = r.db.QueryRowx(query, args...).Err()
	if err != nil {
		return nil, jsend.NewError("failed QueryRowx", err, http.StatusInternalServerError)
	}
	return account, nil
}
func (r *accountRepository) Delete(id int, userID int) error {
	query, args, err := sq.Delete(tableAccounts).
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
