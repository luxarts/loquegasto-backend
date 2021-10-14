package defines

const (
	// Transactions
	EndpointTransactionsCreate         = "/transaction"
	EndpointTransactionsUpdateByMsgID  = "/transactions/:" + ParamMsgID
	EndpointTransactionsGetAllByUserID = "/transactions"

	EndpointPing = "/ping"
)
