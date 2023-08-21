package repository

import (
	"fmt"
	"loquegasto-backend/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestTransactionsSQL_CreateSQL(t *testing.T) {
	// Given
	var catID int64 = 1
	createdAt := time.Now()
	transaction := domain.Transaction{
		ID:          "uuid",
		UserID:      123,
		MsgID:       456,
		Amount:      123456789,
		Description: "test",
		WalletID:    1,
		CreatedAt:   &createdAt,
		CategoryID:  &catID,
	}
	expectedValues := fmt.Sprintf("[%v %v %v %v %v %v %v %v]",
		transaction.ID,
		transaction.UserID,
		transaction.MsgID,
		transaction.Amount,
		transaction.Description,
		transaction.WalletID,
		transaction.CreatedAt,
		transaction.CategoryID)

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.CreateSQL(&transaction)

	// Then
	require.Nil(t, err)
	assert.Equal(t, "INSERT INTO core.transactions (uuid,user_id,msg_id,amount,description,wallet_id,created_at,category_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", query)
	assert.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
func TestTransactionsSQL_UpdateByMsgIDSQL(t *testing.T) {
	// Given
	var catID int64 = 1
	createdAt := time.Now()
	transaction := domain.Transaction{
		ID:          "uuid",
		UserID:      123,
		MsgID:       456,
		Amount:      123456789,
		Description: "test",
		WalletID:    1,
		CreatedAt:   &createdAt,
		CategoryID:  &catID,
	}
	expectedValues := fmt.Sprintf("[%v %v %v %v %v %v %v %v %v %v]",
		transaction.ID,
		transaction.UserID,
		transaction.MsgID,
		transaction.Amount,
		transaction.Description,
		transaction.WalletID,
		transaction.CreatedAt,
		transaction.CategoryID,
		transaction.MsgID,
		transaction.UserID)

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.UpdateByMsgIDSQL(&transaction)

	// Then
	require.Nil(t, err)
	require.Equal(t, "UPDATE core.transactions SET uuid = $1, user_id = $2, msg_id = $3, amount = $4, description = $5, wallet_id = $6, created_at = $7, category_id = $8 WHERE (msg_id = $9 AND user_id = $10)", query)
	require.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
func TestTransactionsSQL_GetAllSQL(t *testing.T) {
	// Given
	var userID int64 = 123
	var msgID int64 = 456
	expectedValues := fmt.Sprintf("[%v %v]", msgID, userID)

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.GetByMsgIDSQL(msgID, userID)

	// Then
	require.Nil(t, err)
	require.Equal(t, "SELECT * FROM core.transactions WHERE (msg_id = $1 AND user_id = $2)", query)
	require.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
func TestTransactionsSQL_GetByMsgIDSQL(t *testing.T) {
	// Given
	var userID int64 = 123
	filters := domain.TransactionFilters{
		"category_id": "1",
		"wallet_id":   "2",
	}
	expectedValues := fmt.Sprintf("[%v %v %v]", userID, "1", "2")

	tsql := transactionsSQL{}
	// When
	query, args, err := tsql.GetAllSQL(userID, &filters)

	// Then
	require.Nil(t, err)
	require.Equal(t, "SELECT * FROM core.transactions WHERE (user_id = $1 AND category_id = $2 AND wallet_id = $3) ORDER BY created_at DESC", query)
	require.Equal(t, expectedValues, fmt.Sprintf("%v", args))
}
