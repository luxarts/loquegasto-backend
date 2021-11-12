package repository

import (
	"fmt"
	"loquegasto-backend/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransactionsSQL_CreateSQL(t *testing.T) {
	// Given
	createdAt := time.Now()
	transaction := domain.Transaction{
		ID:          "uuid",
		UserID:      123,
		MsgID:       456,
		Amount:      123456789,
		Description: "test",
		WalletID:    1,
		CreatedAt:   &createdAt,
	}
	expectedValues := fmt.Sprintf("[%s %d %d %d %s %d %s]",
		transaction.ID,
		transaction.UserID,
		transaction.MsgID,
		transaction.Amount,
		transaction.Description,
		transaction.WalletID,
		transaction.CreatedAt)

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.CreateSQL(&transaction)

	// Then
	require.Nil(t, err)
	require.Equal(t, "INSERT INTO backend.transactions (uuid,user_id,msg_id,amount,description,wallet_id,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)", query)
	require.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
func TestTransactionsSQL_UpdateByMsgIDSQL(t *testing.T) {
	// Given
	createdAt := time.Now()
	transaction := domain.Transaction{
		ID:          "uuid",
		UserID:      123,
		MsgID:       456,
		Amount:      123456789,
		Description: "test",
		WalletID:    1,
		CreatedAt:   &createdAt,
	}
	expectedValues := fmt.Sprintf("[%s %d %d %d %s %d %s]",
		transaction.ID,
		transaction.UserID,
		transaction.MsgID,
		transaction.Amount,
		transaction.Description,
		transaction.WalletID,
		transaction.CreatedAt)

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.CreateSQL(&transaction)

	// Then
	require.Nil(t, err)
	require.Equal(t, "INSERT INTO backend.transactions (uuid,user_id,msg_id,amount,description,wallet_id,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)", query)
	require.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
