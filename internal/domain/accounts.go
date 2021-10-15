package domain

import (
	"time"
)

type Account struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	Name      string     `db:"name"`
	Balance   int64      `db:"balance"`
	CreatedAt *time.Time `db:"created_at"`
}
type AccountDTO struct {
	ID        int        `json:"id"`
	UserID    int        `json:"userID"`
	Name      string     `json:"name,omitempty"`
	Balance   float64    `json:"balance"`
	CreatedAt *time.Time `json:"updated_at"`
}

func (a *Account) ToDTO() *AccountDTO {
	return &AccountDTO{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Balance:   float64(a.Balance) / 100.0,
		CreatedAt: a.CreatedAt,
	}
}
func (dto *AccountDTO) ToAccount() *Account {
	return &Account{
		ID:        dto.ID,
		UserID:    dto.UserID,
		Name:      dto.Name,
		Balance:   int64(dto.Balance * 100),
		CreatedAt: dto.CreatedAt,
	}
}
func (dto *AccountDTO) IsValid() bool {
	return dto.UserID != 0 &&
		dto.Name != "" &&
		dto.CreatedAt != nil
}
