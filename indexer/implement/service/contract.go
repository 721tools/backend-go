package service

import (
	"context"

	"github.com/ethereum/go-ethereum/common/math"

	model2 "github.com/721tools/backend-go/indexer/implement/model"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/model"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"xorm.io/xorm"
)

type Contract struct {
	xorm   *xorm.Engine
	client client.Client
}

func NewContract() *Contract {
	return &Contract{xorm: db.GetDBEngine(), client: client.GetClient()}
}

func (c *Contract) GetContractType(ctx context.Context, address hex.Hex) model.ContractType {
	var contractInfo model2.Contract
	get, _ := c.xorm.Table(&model2.Contract{}).Where("address = ?", address).Get(&contractInfo)
	if get && contractInfo.Address != nil {
		return model.NewContractType(contractInfo.Type)
	}
	return model.EOA
}

func (b *BlockService) Done(ctx context.Context, blockHeights []math.HexOrDecimal64) error {
	var heights []uint64
	for _, height := range blockHeights {
		heights = append(heights, uint64(height))
	}
	_, err := b.xorm.Table(&model2.OriginBlock{}).In("height", heights).Update(map[string]interface{}{
		"done": true,
	})
	return err
}

func (c *Contract) Sync(ctx context.Context, contract *model2.Contract) error {
	var contractInfo model2.Contract
	get, _ := c.xorm.Table(&model2.Contract{}).Where("address = ?", contract.Address).Get(&contractInfo)
	if get && contractInfo.Address != nil {
		return nil
	}

	rawContractInfo, err := c.client.GetContract(ctx, contract.Address)
	if err != nil {
		log.Err(ctx, "get contract occur a error", "contract addrs", contract.Address, "err", err)
	}
	log.Info("contract info", "info", rawContractInfo)
	if rawContractInfo.Type.IsEOA() {
		return nil
	}

	if len(rawContractInfo.Name) > 256 || len(rawContractInfo.Symbol) > 256 {
		log.Warn("contract info", "name or symbol tooo long", len(rawContractInfo.Name))
		return nil
	}

	contractInfo.Height = contract.Height
	contractInfo.TxHash = contract.TxHash
	contractInfo.Address = contract.Address
	contractInfo.Type = rawContractInfo.Type.ToUInt()
	contractInfo.IsStandard = true
	contractInfo.Name = rawContractInfo.Name
	contractInfo.Symbol = rawContractInfo.Symbol
	contractInfo.TotalSupply = hex.BigInt(rawContractInfo.TotalSupply)
	contractInfo.Decimal = uint64(rawContractInfo.Decimals)

	session := c.xorm.NewSession()
	if err = session.Begin(); err != nil {
		log.Err(ctx, "begin contract session failed", "address", contract, "err", err)
		return err
	}
	defer session.Close()
	if _, err = session.Table(&model2.Contract{}).Insert(&contractInfo); err != nil {
		log.Err(ctx, "insert contract data to db failed", "address", contract, "err", err)
		return err
	}
	if err = session.Commit(); err != nil {
		log.Err(ctx, "commit contract data to db failed", "address", contract, "err", err)
		return err
	}
	return nil
}
