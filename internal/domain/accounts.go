package domain

import (
	"time"
)

type Account struct {
	ID        int
	UserID    int
	Name      string
	Balance   int64
	UpdatedAt *time.Time
}
type AccountDTO struct {
	ID        int        `json:"id"`
	UserID    int        `json:"userID"`
	Name      string     `json:"name"`
	Balance   float64    `json:"balance"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (a *Account) ToDTO() *AccountDTO {
	return &AccountDTO{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Balance:   float64(a.Balance) / 100.0,
		UpdatedAt: a.UpdatedAt,
	}
}
func (dto *AccountDTO) ToAccount() *Account {
	return &Account{
		ID:        dto.ID,
		UserID:    dto.UserID,
		Name:      dto.Name,
		Balance:   int64(dto.Balance * 100),
		UpdatedAt: dto.UpdatedAt,
	}
}
func (dto *AccountDTO) IsValid() bool {
	return dto.UserID != 0 &&
		dto.Name != "" &&
		dto.UpdatedAt != nil
}
