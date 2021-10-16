package domain

import "time"

type User struct {
	ID        int
	CreatedAt *time.Time
	ChatID    int
}
type UserDTO struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	ChatID    int        `json:"chat_id"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		ChatID:    u.ChatID,
	}
}
func (dto *UserDTO) ToUser() *User {
	return &User{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		ChatID:    dto.ChatID,
	}
}

func (dto *UserDTO) IsValid() bool {
	return dto.ID != 0 &&
		dto.ChatID != 0 &&
		dto.CreatedAt != nil
}
