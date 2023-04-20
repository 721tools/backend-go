package pulling

import (
	"context"
	"testing"

	client2 "github.com/721tools/backend-go/indexer/pkg/blockchain/client"
	"github.com/stretchr/testify/assert"
)

func TestBlockData_GetBlocksFromNode(t *testing.T) {
	client := client2.NewEthRpcClient("https://cloudflare-eth.com")
	data := NewBlockData(context.Background(), 1, 3, true, client)
	err := data.GetBlocksFromNode()
	assert.NoError(t, err)
	t.Log(data.Data)
}
