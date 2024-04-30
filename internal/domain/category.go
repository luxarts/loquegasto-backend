package domain

import (
	"loquegasto-backend/internal/utils/sanitizer"
	"time"
)

type Category struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	Name          string    `db:"name"`
	SanitizedName string    `db:"sanitized_name"`
	Emoji         string    `db:"emoji"`
	CreatedAt     time.Time `db:"created_at"`
	Deleted       *bool     `db:"deleted"`
}

type CategoryCreateRequest struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}
type CategoryCreateResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	CreatedAt     time.Time `json:"created_at"`
}

func (c *Category) ToCategoryCreateResponse() *CategoryCreateResponse {
	return &CategoryCreateResponse{
		ID:            c.ID,
		Name:          c.Name,
		SanitizedName: c.SanitizedName,
		Emoji:         c.Emoji,
		CreatedAt:     c.CreatedAt,
	}
}

func (req *CategoryCreateRequest) IsValid() bool {
	return req.Name != ""
}

func (req *CategoryCreateRequest) ToCategory() *Category {
	return &Category{
		Name:  req.Name,
		Emoji: req.Emoji,
	}
}

type CategoryUpdateRequest struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}
type CategoryUpdateResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	CreatedAt     time.Time `json:"created_at"`
}

func (req *CategoryUpdateRequest) ToCategory() *Category {
	return &Category{
		Name:          req.Name,
		Emoji:         req.Emoji,
		SanitizedName: sanitizer.Sanitize(req.Name),
	}
}
func (c *Category) ToCategoryUpdateResponse() *CategoryUpdateResponse {
	return &CategoryUpdateResponse{
		ID:            c.ID,
		Name:          c.Name,
		SanitizedName: c.SanitizedName,
		Emoji:         c.Emoji,
		CreatedAt:     c.CreatedAt,
	}
}
func (req *CategoryUpdateRequest) IsValid() bool {
	return req.Name != ""
}
