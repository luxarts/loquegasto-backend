package domain

import (
	"time"
)

type Wallet struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	Name          string    `db:"name"`
	SanitizedName string    `db:"sanitized_name"`
	Emoji         string    `db:"emoji"`
	Balance       int64     `db:"balance"`
	CreatedAt     time.Time `db:"created_at"`
}
type WalletCreateRequest struct {
	Name          string  `json:"name"`
	InitialAmount float64 `json:"balance"`
	Emoji         string  `json:"emoji"`
}

type WalletCreateResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	SanitizedName string    `json:"sanitized_name"`
	Emoji         string    `json:"emoji"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

func (w *Wallet) ToWalletCreateResponse() *WalletCreateResponse {
	return &WalletCreateResponse{
		ID:            w.ID,
		Name:          w.Name,
		Emoji:         w.Emoji,
		SanitizedName: w.SanitizedName,
		Balance:       float64(w.Balance) / 100.0,
		CreatedAt:     w.CreatedAt,
	}
}
func (req *WalletCreateRequest) ToWallet() *Wallet {
	return &Wallet{
		Name:    req.Name,
		Balance: int64(req.InitialAmount * 100),
		Emoji:   req.Emoji,
	}
}
func (req *WalletCreateRequest) IsValid() bool {
	return req.Name != ""
}
