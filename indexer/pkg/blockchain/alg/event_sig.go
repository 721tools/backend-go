package alg

var (
	TransferEventSigID = EventSig("Transfer(address,address,uint256)")
	ApprovalEventSigID = EventSig("Approval(address,address,uint256)")

	OrderFulfilledSigId = EventSig("OrderFulfilled(bytes32,address,address,address,(uint8,address,uint256,uint256)[],(uint8,address,uint256,uint256,address)[])")
	OrdersMatchedSigId  = EventSig("OrdersMatched(address,address,(address,uint8,address,address,uint256,uint256,address,uint256,uint256,uint256,(uint16,address)[],uint256,bytes),bytes32,(address,uint8,address,address,uint256,uint256,address,uint256,uint256,uint256,(uint16,address)[],uint256,bytes),bytes32)")
)
