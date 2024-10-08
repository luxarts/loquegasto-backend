package repository

import (
	"database/sql"
	"errors"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"

	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"
)

const (
	tableUsers = "core.users"
)

type UsersRepository interface {
	Create(user *domain.User) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	GetByChatID(id int64) (*domain.User, error)
}
type usersRepository struct {
	db         *sqlx.DB
	sqlBuilder *usersSQL
}

func NewUsersRepository(db *sqlx.DB) UsersRepository {
	return &usersRepository{
		db:         db,
		sqlBuilder: &usersSQL{},
	}
}

func (r *usersRepository) Create(u *domain.User) (*domain.User, error) {
	query, args, err := r.sqlBuilder.CreateSQL(u)
	if err != nil {
		return nil, jsend.NewError("failed usersRepository.Create.CreateSQL", err, http.StatusInternalServerError)
	}
	_, err = r.db.Exec(query, args...)
	if err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			if pgerr.Code == defines.PGCodeDuplicateKey {
				return nil, jsend.NewError("user ID already exists", nil, http.StatusConflict)
			}
		}
		return nil, jsend.NewError("failed usersRepository.Create.Exec", err, http.StatusInternalServerError)
	}

	return u, nil
}
func (r *usersRepository) GetByID(id string) (*domain.User, error) {
	query, args, err := r.sqlBuilder.GetByIDSQL(id)
	if err != nil {
		return nil, jsend.NewError("failed usersRepository.GetByID.GetByIDSQL", err, http.StatusInternalServerError)
	}

	var user domain.User
	err = r.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, jsend.NewError("user not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed usersRepository.GetByID.StructScan", err, http.StatusInternalServerError)
	}

	return &user, nil
}
func (r *usersRepository) GetByChatID(id int64) (*domain.User, error) {
	query, args, err := r.sqlBuilder.GetByChatIDSQL(id)
	if err != nil {
		return nil, jsend.NewError("failed usersRepository.GetByChatID.GetByChatIDSQL", err, http.StatusInternalServerError)
	}

	var user domain.User
	err = r.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, jsend.NewError("user not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed usersRepository.GetByChatID.StructScan", err, http.StatusInternalServerError)
	}

	return &user, nil
}

// SQL builders
type usersSQL struct{}

func (usql *usersSQL) CreateSQL(u *domain.User) (string, []interface{}, error) {
	return sq.Insert(tableUsers).
		Columns(dbstruct.GetColumns(u)...).
		Values(dbstruct.GetValues(u)...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (usql *usersSQL) GetByIDSQL(id string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableUsers).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (usql *usersSQL) GetByChatIDSQL(id int64) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableUsers).
		Where(sq.Eq{"chat_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
