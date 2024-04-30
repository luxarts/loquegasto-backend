package repository

import (
	"errors"
	"loquegasto-backend/internal/defines"
	"loquegasto-backend/internal/domain"
	"loquegasto-backend/internal/utils/dbstruct"
	"net/http"
	"strings"

	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"

	"github.com/luxarts/jsend-go"

	"github.com/jmoiron/sqlx"
)

const (
	tableCategories = "core.categories"
)

type CategoriesRepository interface {
	Create(category *domain.Category) (*domain.Category, error)
	GetAll(userID string) (*[]domain.Category, error)
	GetByID(id string, userID string) (*domain.Category, error)
	GetByName(name string, userID string) (*domain.Category, error)
	GetBySanitizedName(name string, userID string) (*domain.Category, error)
	GetByEmoji(emoji string, userID string) (*domain.Category, error)
	DeleteByID(id string, userID string) error
	UpdateByID(category *domain.Category, id string, userID string) (*domain.Category, error)
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
		return nil, jsend.NewError("failed categoriesRepository.Create.CreateSQL", err, http.StatusInternalServerError)
	}
	_, err = r.db.Exec(query, args...)
	if err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			if pgerr.Code == defines.PGCodeDuplicateKey {
				return nil, jsend.NewError("wallet ID already exists", nil, http.StatusConflict)
			}
		}
		return nil, jsend.NewError("failed categoriesRepository.Create.Exec", err, http.StatusInternalServerError)
	}

	return category, nil
}
func (r *categoriesRepository) GetAll(userID string) (*[]domain.Category, error) {
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
func (r *categoriesRepository) GetByID(id string, userID string) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.GetByIDSQL(id, userID)
	if err != nil {
		return nil, jsend.NewError("failed categoriesRepository.GetByID.GetByIDSQL", err, http.StatusInternalServerError)
	}

	var category domain.Category
	err = r.db.QueryRowx(query, args...).StructScan(&category)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("category not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed categoriesRepository.GetByID.StructScan", err, http.StatusInternalServerError)
	}

	return &category, nil
}
func (r *categoriesRepository) GetByName(name string, userID string) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.GetByNameSQL(name, userID)
	if err != nil {
		return nil, jsend.NewError("failed GetByNameSQL", err, http.StatusInternalServerError)
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
func (r *categoriesRepository) GetBySanitizedName(name string, userID string) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.GetBySanitizedNameSQL(name, userID)

	if err != nil {
		return nil, jsend.NewError("failed categoriesRepository.GetBySanitizedName.GetBySanitizedNameSQL", err, http.StatusInternalServerError)
	}

	var category domain.Category
	err = r.db.QueryRowx(query, args...).StructScan(&category)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, jsend.NewError("category not found", nil, http.StatusNotFound)
		}
		return nil, jsend.NewError("failed categoriesRepository.GetBySanitizedName.StructScan", err, http.StatusInternalServerError)
	}

	return &category, nil
}
func (r *categoriesRepository) GetByEmoji(emoji string, userID string) (*domain.Category, error) {
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
func (r *categoriesRepository) DeleteByID(id string, userID string) error {
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
func (r *categoriesRepository) UpdateByID(category *domain.Category, id string, userID string) (*domain.Category, error) {
	query, args, err := r.sqlBuilder.UpdateByIDSQL(category, id, userID)

	if err != nil {
		return nil, jsend.NewError("failed UpdateSQL", err, http.StatusInternalServerError)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, jsend.NewError("failed Exec", err, http.StatusInternalServerError)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return nil, jsend.NewError("category not found", nil, http.StatusNotFound)
	}

	category.ID = id
	return category, nil
}

// SQL Builders
type categoriesSQL struct{}

func (csql *categoriesSQL) CreateSQL(category *domain.Category) (string, []interface{}, error) {
	return sq.Insert(tableCategories).
		Columns(dbstruct.GetColumns(category)...).
		Values(dbstruct.GetValues(category)...).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetAllSQL(userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.Eq{"user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetByIDSQL(id string, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Eq{"user_id": userID, "deleted": nil},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetByNameSQL(name string, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.Eq{"sanitized_name": name, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetBySanitizedNameSQL(name string, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.Eq{"sanitized_name": name, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) GetByEmojiSQL(emoji string, userID string) (string, []interface{}, error) {
	return sq.Select("*").
		From(tableCategories).
		Where(sq.Eq{"emoji": emoji, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) DeleteByIDSQL(id string, userID string) (string, []interface{}, error) {
	return sq.Update(tableCategories).
		Set("deleted", true).
		Where(sq.Eq{"id": id, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
func (csql *categoriesSQL) UpdateByIDSQL(category *domain.Category, id string, userID string) (string, []interface{}, error) {
	return sq.Update(tableCategories).
		Set("name", category.Name).
		Set("sanitized_name", category.SanitizedName).
		Set("emoji", category.Emoji).
		Where(sq.Eq{"id": id, "user_id": userID, "deleted": nil}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
}
