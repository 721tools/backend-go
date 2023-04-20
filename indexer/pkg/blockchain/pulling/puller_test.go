package pulling

import (
	"testing"
	"time"

	client2 "github.com/721tools/backend-go/index/pkg/blockchain/client"
	"github.com/stretchr/testify/assert"
)

func TestPuller_Continuous(t *testing.T) {
	client := client2.NewEthRpcClient("https://cloudflare-eth.com")
	puller := NewPuller(3, 30, client)
	go puller.Continuous(3)
	for {
		pop := puller.Pop()
		for idx := range pop.Data {
			t.Log("block height", pop.Data[idx].Number)
		}
		time.Sleep(time.Second)
	}
}

func TestPuller_Temporarily(t *testing.T) {
	client := client2.NewEthRpcClient("https://cloudflare-eth.com")
	puller := NewPuller(3, 30, client)
	puller.Temporarily(3, 5)
	for {
		pop := puller.Pop()
		if pop.Cancel {
			break
		}
		data := pop.Data
		for idx := range data {
			t.Log("idx", idx, "height", data[idx].Number)
		}
		assert.Equal(t, 2, len(data))
	}
}
