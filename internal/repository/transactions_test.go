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
func TestTransactionsSQL_GetAllByUserIDSQL(t *testing.T) {
	// Given
	userID := 123
	filters := domain.TransactionFilters{
		"category_id": "1",
		"wallet_id":   "2",
	}
	expectedValues := fmt.Sprintf("[%d %s %s]", userID, "1", "2")

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.GetAllByUserIDSQL(userID, &filters)

	// Then
	require.Nil(t, err)
	require.Equal(t, "SELECT * FROM backend.transactions WHERE (user_id = $1 AND category_id = $2 AND wallet_id = $3)", query)
	require.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
