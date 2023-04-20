package block_listener

import (
	db2 "github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
)

var log = log16.NewLogger("module", "block_listener")

type OriginDBStartListener struct {
}

func NewOriginDBStartListener() *OriginDBStartListener {
	return &OriginDBStartListener{}
}

// Handle
// deal with origin_block, origin_tx, origin_receipt_log
func (db *OriginDBStartListener) Handle(event *Event) (err error) {
	session := db2.GetDBEngine().NewSession()
	if err = session.Begin(); err != nil {
		return err
	}
	// block
	if _, err = session.Insert(event.Block); err != nil {
		return err
	}

	// txs
	if len(event.Block.Txs) > 0 {
		if _, err = session.InsertMulti(&event.Block.Txs); err != nil {
			return err
		}
	}

	// event logs
	for _, tx := range event.Block.Txs {
		if len(tx.ReceiptLogs) == 0 {
			continue
		}
		explode := explodeG(tx.ReceiptLogs, receiptLogsChunkSize)
		for _, arr := range explode {
			if _, err = session.InsertMulti(&arr); err != nil {
				return err
			}
		}
	}
	return session.Commit()
}
