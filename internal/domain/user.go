package domain

import "time"

type User struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ChatID    int       `db:"chat_id"`
	State     *string   `db:"state"`
}
type UserDTO struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ChatID    int        `json:"chat_id"`
	State     *string    `json:"state"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		CreatedAt: &u.CreatedAt,
		UpdatedAt: &u.UpdatedAt,
		ChatID:    u.ChatID,
		State:     u.State,
	}
}
func (dto *UserDTO) ToUser() *User {
	return &User{
		ID:        dto.ID,
		CreatedAt: *dto.CreatedAt,
		UpdatedAt: *dto.UpdatedAt,
		ChatID:    dto.ChatID,
		State:     dto.State,
	}
}

func (dto *UserDTO) IsValid() bool {
	return dto.ID != 0 &&
		dto.ChatID != 0 &&
		dto.CreatedAt != nil &&
		dto.UpdatedAt != nil
}
