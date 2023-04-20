package abi_parse

import (
	"github.com/721tools/backend-go/indexer/pkg/blockchain/alg"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

type ERC20 struct {
}

func NewERC20() *ERC20 {
	return &ERC20{}
}

func (e *ERC20) Tag() string {
	return "ERC20"
}

func (e *ERC20) Address() []hex.Hex {
	return make([]hex.Hex, 0)
}

func (e *ERC20) Methods() map[string]JieHandle {
	return map[string]JieHandle{
		alg.TransferMethod.ToString():     e.transferMethod,
		alg.TransferFromMethod.ToString(): e.transferFromMethod,
		alg.ApproveMethod.ToString():      e.approveMethod,
	}
}

func (e *ERC20) EventLogs() map[string]JieHandle {
	return map[string]JieHandle{
		alg.TransferEventSigID.ToString(): e.transferEvent,
	}
}

func (e *ERC20) transferMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 2 {
		return "", nil
	}
	args := make(Args)
	args["to"] = params[0].ToAddress()
	args["value"] = params[1].ToHexInt()
	return "transfer", args
}

func (e *ERC20) transferFromMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["from"] = params[0].ToAddress()
	args["to"] = params[1].ToAddress()
	args["amount"] = params[2].ToHexInt()
	return "transferFrom", args
}

func (e *ERC20) approveMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 2 {
		return "", nil
	}
	args := make(Args)
	args["spender"] = params[0].ToAddress()
	args["amount"] = params[1].ToHexInt()
	return "approve", args
}

func (e *ERC20) transferEvent(hexData hex.Hex) (string, Args) {
	params := dealWithEventLogInput(hexData)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["from"] = params[0].ToAddress()
	args["to"] = params[1].ToAddress()
	args["value"] = params[2].ToHexInt()
	args["amount"] = params[2].ToHexInt()
	return "Transfer", args
}
