package repository

import (
	"database/sql"
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
	GetByID(id int) (*domain.User, error)
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
	_, err = r.db.Exec(query, args...)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == defines.PGCodeDuplicateKey {
				return nil, jsend.NewError("user ID already exists", nil, http.StatusConflict)
			}
		}
		return nil, jsend.NewError("failed CreateSQL", err, http.StatusInternalServerError)
	}

	return u, nil
}
func (r *usersRepository) GetByID(id int) (*domain.User, error) {
	query, args, err := r.sqlBuilder.GetByIDSQL(id)
	if err != nil {
		return nil, jsend.NewError("failed GetByIDSQL", err, http.StatusInternalServerError)
	}

	var user domain.User
	err = r.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, jsend.NewError("user not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
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
func (usql *usersSQL) GetByIDSQL(id int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableUsers).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
