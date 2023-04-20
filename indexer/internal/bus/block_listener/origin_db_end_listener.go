package block_listener

import (
	model2 "github.com/721tools/backend-go/indexer/internal/model"
	db2 "github.com/721tools/backend-go/indexer/pkg/db"
)

type OriginDBEndListener struct {
}

func NewOriginDBEndListener() *OriginDBEndListener {
	return &OriginDBEndListener{}
}

// Handle
// deal with origin_block, origin_tx, origin_receipt_log
func (db *OriginDBEndListener) Handle(event *Event) (err error) {
	session := db2.NewSession()
	if err = session.Begin(); err != nil {
		log.Error("begin session occur a err", "err", err)
		return err
	}
	_, err = session.Table(&model2.OriginBlock{}).Where("height = ?", event.Block.Height).Update(map[string]interface{}{
		"done": true,
	})
	if err != nil {
		log.Error("begin session occur a err", "err", err)
		return err
	}
	return session.Commit()
}
