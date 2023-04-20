package block_listener

import (
	model2 "github.com/721tools/backend-go/indexer/internal/model"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
	"github.com/ethereum/go-ethereum/common/math"
)

const (
	receiptLogsChunkSize      chunkSize = 50
	tokenFlowsChunkSize       chunkSize = 100
	nftMintChunkSize          chunkSize = 100
	nftSalesChunkSize         chunkSize = 100
	contractCreationChunkSize chunkSize = 100
)

type explodeType interface {
	model2.OriginReceiptLog | model2.TokenFlow | model2.NFTMint | model2.NFTSales | model2.Contract
}

func explodeG[T explodeType](arr []T, num chunkSize) [][]T {
	var segments = make([][]T, 0)
	max := int64(len(arr))
	num = int64(num)

	if max == 0 {
		return segments
	}

	if max < num {
		segments = append(segments, arr)
		return segments
	}

	quantity := max / num
	for i := int64(0); i <= quantity; i++ {

		start := i * num
		end := (i + 1) * num
		if start == max {
			break
		}
		if end > max {
			end = max
		}
		segments = append(segments, arr[start:end])
	}
	return segments
}

func updateLog(txHash hex.Hex, logIndex math.HexOrDecimal64, name []byte, args []byte) {
	_, _ = db.GetDBEngine().Table(&model2.OriginReceiptLog{}).Where("tx_hash = ? AND log_index = ?", txHash, uint64(logIndex)).Limit(1).Update(map[string]interface{}{
		"event_name":   name,
		"EventPayload": args,
	})
}
