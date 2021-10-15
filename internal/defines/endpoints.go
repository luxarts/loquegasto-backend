package defines

const (
	EndpointPing = "/ping"

	// Transactions
	EndpointTransactionsCreate         = "/transaction"
	EndpointTransactionsUpdateByMsgID  = "/transactions/:" + ParamMsgID
	EndpointTransactionsGetAllByUserID = "/transactions"

	// Users
	EndpointUsersCreate = "/user"

	// Account
	EndpointAccountsCreate     = "/account"
	EndpointAccountsGetByName  = "/accounts"
	EndpointAccountsGetByID    = "/accounts/:id"
	EndpointAccountsUpdateByID = "/accounts/:id"
	EndpointAccountsDeleteByID = "/accounts/:id"
)
