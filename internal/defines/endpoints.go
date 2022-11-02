package defines

const (
	EndpointPing = "/ping"

	// Authorize
	EndpointAuthorize = "/auth/:" + ParamUserID

	// Transactions
	EndpointTransactionsCreate        = "/transactions"
	EndpointTransactionsUpdateByMsgID = "/transaction/:" + ParamMsgID
	EndpointTransactionsGetAll        = "/transactions"

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
