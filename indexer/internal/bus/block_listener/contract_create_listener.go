package block_listener

import (
	"context"

	model2 "github.com/721tools/backend-go/indexer/internal/model"
	"github.com/721tools/backend-go/indexer/internal/service"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"xorm.io/xorm"
)

type ContractCreateListener struct {
	svc  service.ContractIface
	xorm *xorm.Engine
}

func NewContractCreateListener() *ContractCreateListener {
	return &ContractCreateListener{svc: service.NewContract(), xorm: db.GetDBEngine()}
}

func (t *ContractCreateListener) Handle(event *Event) (err error) {
	for _, rawTx := range event.Block.Txs {
		if rawTx.RawTxReceipt.ContractAddress != nil {
			contract := model2.Contract{
				Height:  rawTx.Height,
				TxHash:  rawTx.TxHash,
				Address: rawTx.RawTxReceipt.ContractAddress,
			}
			return t.svc.Sync(context.Background(), &contract)
		}
	}
	return nil
}
