package model

import (
	"github.com/721tools/backend-go/index/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
)

const EOA = ContractType(0)
const ERC20 = ContractType(1)
const ERC721 = ContractType(2)
const ERC1155 = ContractType(3)
const Contract = ContractType(4)

type ContractType uint

func NewContractType(t uint) ContractType {
	return ContractType(t)
}

func (c ContractType) IsERC721() bool {
	return c == ERC721
}

func (c ContractType) IsERC20() bool {
	return c == ERC20
}

func (c ContractType) IsEOA() bool {
	return c == EOA
}

func (c ContractType) IsERC1155() bool {
	return c == ERC1155
}

func (c ContractType) IsToken() bool {
	return c.IsERC20() || c.IsERC721() || c.IsERC1155()
}

func (c ContractType) ToUInt() uint {
	return uint(c)
}

type RawContract struct {
	Address     hex.Hex
	Type        ContractType
	Name        string
	Symbol      string
	Decimals    math.HexOrDecimal64
	TotalSupply math.HexOrDecimal256
}
