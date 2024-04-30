package domain

import "time"

type User struct {
	ID             string    `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	ChatID         int64     `db:"chat_id"`
	TimezoneOffset int       `db:"timezone_offset"`
}
type UserCreateRequest struct {
	ChatID         int64 `json:"chat_id"`
	TimezoneOffset int   `json:"timezone_offset"`
}
type UserCreateResponse struct {
	ID             string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	ChatID         int64     `json:"chat_id"`
	TimezoneOffset int       `json:"timezone_offset"`
}

func (u *User) ToResponse() *UserCreateResponse {
	return &UserCreateResponse{
		ID:             u.ID,
		CreatedAt:      u.CreatedAt,
		ChatID:         u.ChatID,
		TimezoneOffset: u.TimezoneOffset,
	}
}
func (body *UserCreateRequest) ToUser() *User {
	return &User{
		ChatID:         body.ChatID,
		TimezoneOffset: body.TimezoneOffset,
	}
}

func (body *UserCreateRequest) IsValid() bool {
	return body.ChatID != 0
}

type UserAuthWithTelegramRequest struct {
	ChatID int64 `json:"chat_id"`
}

func (u *UserAuthWithTelegramRequest) IsValid() bool {
	return u.ChatID != 0
}

type UserAuthWithTelegramResponse struct {
	AccessToken string `json:"access_token"`
}
