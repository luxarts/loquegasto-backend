package domain

import "time"

type User struct {
	ID           int64     `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	ChatID       int       `db:"chat_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	Expiry       time.Time `db:"expiry"`
}
type UserDTO struct {
	ID           int64      `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	ChatID       int        `json:"chat_id"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	Expiry       *time.Time `json:"expiry"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:           u.ID,
		CreatedAt:    &u.CreatedAt,
		ChatID:       u.ChatID,
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		Expiry:       &u.Expiry,
	}
}
func (dto *UserDTO) ToUser() *User {
	return &User{
		ID:           dto.ID,
		CreatedAt:    *dto.CreatedAt,
		ChatID:       dto.ChatID,
		AccessToken:  dto.AccessToken,
		RefreshToken: dto.RefreshToken,
		Expiry:       *dto.Expiry,
	}
}

func (dto *UserDTO) IsValid() bool {
	return dto.ID != 0 &&
		dto.ChatID != 0 &&
		dto.CreatedAt != nil &&
		dto.AccessToken != "" &&
		dto.RefreshToken != "" &&
		dto.Expiry != nil
}
