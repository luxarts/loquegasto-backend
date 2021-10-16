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
	EndpointAccountsGetAll     = "/accounts"
	EndpointAccountsGetByID    = "/accounts/:" + ParamID
	EndpointAccountsUpdateByID = "/accounts/:" + ParamID
	EndpointAccountsDeleteByID = "/accounts/:" + ParamID
)
