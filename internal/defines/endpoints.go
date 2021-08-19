package defines

const (
	// Transactions
	EndpointTransactionsCreate        = "/transaction"
	EndpointTransactionsGetTotal      = "/transactions/total"
	EndpointTransactionsUpdateByMsgID = "/transactions/:" + ParamMsgID

	EndpointPing = "/ping"
)
