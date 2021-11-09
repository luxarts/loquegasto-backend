package repository

import (
	"loquegasto-backend/internal/domain"
	"net/http"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"
)

const (
	tableUsers = "backend.users"
)

type UsersRepository interface {
	Create(user *domain.User) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
}
type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) UsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) Create(user *domain.User) (*domain.User, error) {
	query := sq.Insert(tableUsers).Columns("id", "chat_id", "created_at").
		Values(
			user.ID,
			user.ChatID,
			user.CreatedAt).
		Suffix("RETURNING \"id\"").
		RunWith(r.db).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&user.ID)
	if err != nil {
		return nil, jsend.NewError("failed QueryRow", err, http.StatusInternalServerError)
	}

	return user, nil
}
func (r *usersRepository) GetByID(id int) (*domain.User, error) {
	condition := sq.Eq{"id": id}
	query, args, err := sq.Select("*").
		From(tableUsers).
		Where(condition).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	var user domain.User
	err = r.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("user not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed QueryRow", err, http.StatusInternalServerError)
	}

	return &user, nil
}
