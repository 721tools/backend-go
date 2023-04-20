package abi_parse

import (
	"github.com/721tools/backend-go/index/pkg/blockchain/alg"
	"github.com/721tools/backend-go/index/pkg/utils/hex"
)

type ERC721 struct {
}

func NewERC721() *ERC721 {
	return &ERC721{}
}

func (n *ERC721) Tag() string {
	return "ERC721"
}

func (n *ERC721) Address() []hex.Hex {
	return make([]hex.Hex, 0)
}

func (n *ERC721) Methods() map[string]JieHandle {
	return map[string]JieHandle{
		alg.MethodSig("safeTransferFrom(address,address,uint256)").ToString():       n.safeTransferFromMethod,
		alg.MethodSig("safeTransferFrom(address,address,uint256,bytes)").ToString(): n.safeTransferFromMethodV2,
		alg.MethodSig("transferFrom(address,address,uint256)").ToString():           n.transferFromMethod,
		alg.MethodSig("approve(address,uint256)").ToString():                        n.approveMethod,
		alg.MethodSig("setApprovalForAll(address,bool)").ToString():                 n.setApprovalForAllMethod,
		alg.MethodSig("mint(address,uint256)").ToString():                           n.mintMethod,
	}
}

func (n *ERC721) EventLogs() map[string]JieHandle {
	return map[string]JieHandle{
		alg.EventSig("Transfer(address,address,uint256)").ToString():    n.transferEvent,
		alg.EventSig("Approval(address,address,uint256)").ToString():    n.approvalEvent,
		alg.EventSig("ApprovalForAll(address,address,bool)").ToString(): n.approvalForAllEvent,
	}
}

// safeTransferFromMethod
// function safeTransferFrom(address from, address to, uint256 tokenId)
func (n *ERC721) safeTransferFromMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["from"] = params[0].ToAddress()
	args["to"] = params[1].ToAddress()
	args["tokenId"] = params[2].ToHexInt()
	return "safeTransferFrom", args
}

// mintMethod
// mint(address _owner, uint256 _amount)
func (n *ERC721) mintMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 2 {
		return "", nil
	}
	args := make(Args)
	args["_owner"] = params[0].ToAddress()
	args["_amount"] = params[1].ToHexInt()
	return "mint", args
}

// safeTransferFromMethodV2
// safeTransferFrom(address from, address to, uint256 tokenId, bytes data)
func (n *ERC721) safeTransferFromMethodV2(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 4 {
		return "", nil
	}
	args := make(Args)
	args["from"] = params[0].ToAddress()
	args["to"] = params[1].ToAddress()
	args["tokenId"] = params[2].ToHexInt()
	args["data"] = params[3]
	return "safeTransferFrom", args
}

// approveMethod
// approve(address to, uint256 tokenId)
func (n *ERC721) approveMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 2 {
		return "", nil
	}
	args := make(Args)
	args["to"] = params[0].ToAddress()
	args["tokenId"] = params[1].ToHexInt()
	return "approve", args
}

// transferFromMethod 转账 transferFrom(address from, address to, uint256 tokenId)
// like: https://cn.etherscan.com/tx/0x6e7e61160609a7c143e8fefb669b8af8513aad4651208c7dbcda97bf39c9736f
func (n *ERC721) transferFromMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["from"] = params[0].ToAddress()
	args["to"] = params[1].ToAddress()
	args["tokenId"] = params[2].ToHexInt()
	return "transferFrom", args
}

// setApprovalForAllMethod
// setApprovalForAll(address operator, bool _approved)
func (n *ERC721) setApprovalForAllMethod(hex hex.Hex) (string, Args) {
	params := dealWithTxInput(hex)
	if len(params) < 2 {
		return "", nil
	}
	args := make(Args)
	args["operator"] = params[0].ToAddress()
	args["_approved"] = params[1].ToBool()
	return "setApprovalForAll", args
}

// transferEvent
// erc20 is value
// Transfer(from, to, tokenId)
func (n *ERC721) transferEvent(hexData hex.Hex) (string, Args) {
	params := dealWithEventLogInput(hexData)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["from"] = params[0].ToAddress()
	args["to"] = params[1].ToAddress()
	args["value"] = params[2].ToHexInt()
	args["amount"] = hex.IntToHex(uint64(1))
	return "Transfer", args
}

// approvalEvent
// Approval(address owner,address approved,uint256 tokenId)
func (n *ERC721) approvalEvent(hex hex.Hex) (string, Args) {
	params := dealWithEventLogInput(hex)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["owner"] = params[0].ToAddress()
	args["approved"] = params[1].ToAddress()
	args["tokenId"] = params[2].ToHexInt()
	return "Approval", args
}

// approvalForAllEvent
// ApprovalForAll(address owner, address operator, bool approved)
func (n *ERC721) approvalForAllEvent(hex hex.Hex) (string, Args) {
	params := dealWithEventLogInput(hex)
	if len(params) < 3 {
		return "", nil
	}
	args := make(Args)
	args["owner"] = params[0].ToAddress()
	args["operator"] = params[1].ToAddress()
	args["approved"] = params[2].ToBool()
	return "ApprovalForAll", args
}
