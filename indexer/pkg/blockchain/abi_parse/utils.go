package abi_parse

import (
	"github.com/721tools/backend-go/index/pkg/utils/hex"
)

func dealWithTxInput(input []byte) []hex.EthParamHex {
	if len(input) < 4 {
		return make([]hex.EthParamHex, 0)
	}
	data := input[4:]
	if len(data) < 32 {
		return make([]hex.EthParamHex, 0)
	}
	params := make([]hex.EthParamHex, 0)
	for {
		if len(data) < 32 {
			break
		}
		params = append(params, data[:32])
		data = data[32:]
	}
	return params
}

func dealWithEventLogInput(input []byte) []hex.EthParamHex {
	if len(input) < 32 {
		return make([]hex.EthParamHex, 0)
	}
	params := make([]hex.EthParamHex, 0)
	for {
		if len(input) < 32 {
			break
		}
		params = append(params, input[:32])
		input = input[32:]
	}
	return params
}
