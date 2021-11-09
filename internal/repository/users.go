package repository

import (
	"loquegasto-backend/internal/domain"
	"net/http"
	"strings"
	"time"

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
	query, args, err := r.createSQL(user.ID, user.ChatID, user.CreatedAt)
	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
	}

	return user, nil
}
func (r *usersRepository) GetByID(id int) (*domain.User, error) {
	query, args, err := r.getByIDSQL(id)
	if err != nil {
		return nil, jsend.NewError("failed ToSql", err, http.StatusInternalServerError)
	}

	var user domain.User
	err = r.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("user not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
	}

	return &user, nil
}

func (r *usersRepository) getByIDSQL(id int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableUsers).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (r *usersRepository) createSQL(userID int, chatID int, createdAt *time.Time) (string, []interface{}, error) {
	return sq.Insert(tableUsers).
		Columns("id", "chat_id", "created_at").
		Values(userID, chatID, createdAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
