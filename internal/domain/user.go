package domain

import "time"

type User struct {
	ID             int64     `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	ChatID         int       `db:"chat_id"`
	TimezoneOffset int       `db:"timezone_offset"`
}
type UserDTO struct {
	ID             int64      `json:"id"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	ChatID         int        `json:"chat_id"`
	TimezoneOffset int        `json:"timezone_offset"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:             u.ID,
		CreatedAt:      &u.CreatedAt,
		UpdatedAt:      &u.UpdatedAt,
		ChatID:         u.ChatID,
		TimezoneOffset: u.TimezoneOffset,
	}
}
func (dto *UserDTO) ToUser() *User {
	return &User{
		ID:             dto.ID,
		CreatedAt:      *dto.CreatedAt,
		UpdatedAt:      *dto.UpdatedAt,
		ChatID:         dto.ChatID,
		TimezoneOffset: dto.TimezoneOffset,
	}
}

func (dto *UserDTO) IsValid() bool {
	return dto.ID != 0 &&
		dto.ChatID != 0 &&
		dto.CreatedAt != nil &&
		dto.UpdatedAt != nil
}
