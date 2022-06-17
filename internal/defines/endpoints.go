package defines

const (
	EndpointPing = "/ping"

	// Transactions
	EndpointTransactionsCreate         = "/transactions"
	EndpointTransactionsUpdateByMsgID  = "/transaction/:" + ParamMsgID
	EndpointTransactionsGetAllByUserID = "/transactions"

	// Users
	EndpointUsersCreate = "/users"
	EndpointUsersGet    = "/user"

	// Wallets
	EndpointWalletsCreate     = "/wallets"
	EndpointWalletsGetAll     = "/wallets"
	EndpointWalletsGetByID    = "/wallet/:" + ParamID
	EndpointWalletsUpdateByID = "/wallet/:" + ParamID
	EndpointWalletsDeleteByID = "/wallet/:" + ParamID

	// Categories
	EndpointCategoriesCreate     = "/categories"
	EndpointCategoriesGetAll     = "/categories"
	EndpointCategoriesDeleteByID = "/category/:" + ParamID
	EndpointCategoriesUpdateByID = "/category/:" + ParamID
)
