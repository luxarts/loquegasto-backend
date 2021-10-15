package repository

import (
	"loquegasto-backend/internal/domain"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"
)

const (
	tableUsers = "backend.users"
)

type UsersRepository interface {
	Create(user *domain.User) (*domain.User, error)
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
