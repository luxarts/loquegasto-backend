package repository

import (
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"
)

const (
	tableCategories = "backend.categories"
)

type CategoriesRepository interface {
	Create(category *domain.Category) (*domain.Category, error)
	GetAll(userID int) (*[]domain.Category, error)
	GetByName(name string, userID int) (*domain.Category, error)
	GetByEmoji(emoji string, userID int) (*domain.Category, error)
	DeleteByID(id int, userID int) error
}

type categoriesRepository struct {
	db         *sqlx.DB
	sqlBuilder *categoriesSQL
}

func NewCategoriesRepository(db *sqlx.DB) CategoriesRepository {
	return &categoriesRepository{
		db:         db,
		sqlBuilder: &categoriesSQL{},
	}
}

func (r *categoriesRepository) Create(category *domain.Category) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.CreateSQL(category)
	if err != nil {
		return nil, jsend.NewError("failed CreateSQL", err, http.StatusInternalServerError)
	}
	err = r.db.QueryRowx(query, args...).Scan(&category.ID)
	if err != nil {
		return nil, jsend.NewError("failed Scan", err, http.StatusInternalServerError)
	}

	return category, nil
}
func (r *categoriesRepository) GetAll(userID int) (*[]domain.Category, error) {
	query, args, err := r.sqlBuilder.GetAllSQL(userID)
	if err != nil {
		return nil, jsend.NewError("failed GetAllSQL", err, http.StatusInternalServerError)
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Queryx", err, http.StatusInternalServerError)
	}

	var results []domain.Category
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, jsend.NewError("failed Err: error in row", err, http.StatusInternalServerError)
		}
		var category domain.Category
		if err := rows.StructScan(&category); err != nil {
			return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
		}
		results = append(results, category)
	}

	return &results, nil
}
func (r *categoriesRepository) GetByName(name string, userID int) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.GetByNameSQL(name, userID)
	if err != nil {
		return nil, jsend.NewError("failed GetByIDSQL", err, http.StatusInternalServerError)
	}

	var category domain.Category
	err = r.db.QueryRowx(query, args...).StructScan(&category)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("category not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
	}

	return &category, nil
}
func (r *categoriesRepository) GetByEmoji(emoji string, userID int) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.GetByEmojiSQL(emoji, userID)
	if err != nil {
		return nil, jsend.NewError("failed GetByEmojiSQL", err, http.StatusInternalServerError)
	}

	var category domain.Category
	err = r.db.QueryRowx(query, args...).StructScan(&category)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("category not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed StructScan", err, http.StatusInternalServerError)
	}

	return &category, nil
}
func (r *categoriesRepository) DeleteByID(id int, userID int) error {
	query, args, err := r.sqlBuilder.DeleteByIDSQL(id, userID)
	if err != nil {
		return jsend.NewError("failed DeleteByIDSQL", err, http.StatusInternalServerError)
	}

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return jsend.NewError("category not found", err, http.StatusNotFound)
	}

	return nil
}

// SQL Builders
type categoriesSQL struct{}

func (csql *categoriesSQL) CreateSQL(category *domain.Category) (string, []interface{}, error) {
	category.ID = nil
	return sq.Insert(tableCategories).
		Columns(dbstruct.GetColumns(category)[1:]...).
		Values(dbstruct.GetValues(category)[1:]...).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetAllSQL(userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetByNameSQL(name string, userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.And{
			sq.Eq{"sanitized_name": name},
			sq.Eq{"user_id": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetByEmojiSQL(emoji string, userID int) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.And{
			sq.Eq{"emoji": emoji},
			sq.Eq{"user_id": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) DeleteByIDSQL(id int, userID int) (string, []interface{}, error) {
	return sq.Delete(tableCategories).
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Eq{"user_id": userID},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
