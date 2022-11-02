package domain

import (
	"time"
)

type Wallet struct {
	ID            *int64     `db:"id"`
	UserID        int64      `db:"user_id"`
	Name          string     `db:"name"`
	SanitizedName string     `db:"sanitized_name"`
	Balance       int64      `db:"balance"`
	CreatedAt     *time.Time `db:"created_at"`
}
type WalletDTO struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id"`
	Name          string     `json:"name,omitempty"`
	SanitizedName string     `json:"sanitized_name"`
	Balance       float64    `json:"balance"`
	CreatedAt     *time.Time `json:"created_at"`
}

func (a *Wallet) ToDTO() *WalletDTO {
	return &WalletDTO{
		ID:            *a.ID,
		UserID:        a.UserID,
		Name:          a.Name,
		SanitizedName: a.SanitizedName,
		Balance:       float64(a.Balance) / 100.0,
		CreatedAt:     a.CreatedAt,
	}
}
func (dto *WalletDTO) ToWallet() *Wallet {
	return &Wallet{
		ID:            &dto.ID,
		UserID:        dto.UserID,
		Name:          dto.Name,
		SanitizedName: dto.SanitizedName,
		Balance:       int64(dto.Balance * 100),
		CreatedAt:     dto.CreatedAt,
	}
}
func (dto *WalletDTO) IsValid() bool {
	return dto.UserID != 0 &&
		dto.Name != "" &&
		dto.CreatedAt != nil
}
