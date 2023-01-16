package domain

type Category struct {
	ID            *int64 `db:"id"`
	UserID        int64  `db:"user_id"`
	Name          string `db:"name"`
	SanitizedName string `db:"sanitized_name"`
	Emoji         string `db:"emoji"`
}

func (c *Category) ToDTO() *CategoryDTO {
	return &CategoryDTO{
		ID:            *c.ID,
		UserID:        c.UserID,
		Name:          c.Name,
		SanitizedName: c.SanitizedName,
		Emoji:         c.Emoji,
	}
}

type CategoryDTO struct {
	ID            int64  `json:"id"`
	UserID        int64  `json:"user_id"`
	Name          string `json:"name"`
	SanitizedName string `json:"sanitized_name"`
	Emoji         string `json:"emoji"`
}

func (dto *CategoryDTO) IsValid() bool {
	return dto.Name != ""
}

func (dto *CategoryDTO) ToCategory() *Category {
	return &Category{
		ID:            &dto.ID,
		UserID:        dto.UserID,
		Name:          dto.Name,
		SanitizedName: dto.SanitizedName,
		Emoji:         dto.Emoji,
	}
}
