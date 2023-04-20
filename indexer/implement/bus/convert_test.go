package bus

import (
	"context"
	"testing"

	"github.com/721tools/backend-go/indexer/implement/service"
	"github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/721tools/backend-go/indexer/pkg/db"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
)

func TestConvert_NewConvert(t *testing.T) {
	client := client.NewClient("https://wild-empty-darkness.discover.quiknode.pro/87866012a3b481a1851cbb91c9687c44b2b55eaa/", client.EvmClient)
	defer client.Release()
	cvt := NewConvert()
	assert.NotNil(t, cvt)
}

func TestConvert_DispatchEvents(t *testing.T) {

	ctx := context.Background()

	client := client.NewClient("https://wild-empty-darkness.discover.quiknode.pro/87866012a3b481a1851cbb91c9687c44b2b55eaa/", client.EvmClient)
	defer client.Release()
	db.NewDBEngine("root:@tcp(localhost:3306)/test_dev?parseTime=true")

	srv := service.NewBlockService()
	srv.ClearBlocks(ctx, []uint64{15139880})

	var batch []math.HexOrDecimal64
	batch = append(batch, 15139880)
	blocks, err := client.BlocksByHeight(ctx, batch, true)
	assert.NoError(t, err)
	for idx := range blocks {
		t.Log("height", blocks[idx].Number)
		for txIdx := range blocks[idx].Transactions {
			t.Log("tx", blocks[idx].Transactions[txIdx].Hash, "receipt gas used", blocks[idx].Transactions[txIdx].Receipt.GasUsed)
		}
	}

	cvt := NewConvert()
	cvt.SetRawBlocks(blocks)

	Init()
	cvt.DispatchEvents()
}
