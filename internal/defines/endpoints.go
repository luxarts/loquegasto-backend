package defines

const (
	EndpointPing = "/ping"

	// Transactions
	EndpointTransactionsCreate         = "/transaction"
	EndpointTransactionsUpdateByMsgID  = "/transactions/:" + ParamMsgID
	EndpointTransactionsGetAllByUserID = "/transactions"

	// Users
	EndpointUsersCreate = "/user"

	// Wallets
	EndpointWalletsCreate     = "/wallet"
	EndpointWalletsGetAll     = "/wallets"
	EndpointWalletsGetByID    = "/wallets/:" + ParamID
	EndpointWalletsUpdateByID = "/wallets/:" + ParamID
	EndpointWalletsDeleteByID = "/wallets/:" + ParamID
)
