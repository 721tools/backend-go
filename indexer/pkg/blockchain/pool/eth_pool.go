package pool

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/inconshreveable/log15"
	"github.com/silenceper/pool"
)

var IdleTimeOut = 60
var InitCap = 5
var MaxIdle = 20
var MaxCap = 50
var ChunkSize = 1000

type EthRpcPool struct {
	pool.Pool
	opts []Option
}

// NewEthRPCPool 创建eth连接到节点的连接池
func NewEthRPCPool(endpoint string, opts ...Option) (rpcPool EthRpcPool, err error) {
	poolConfig := &pool.Config{
		InitialCap: InitCap,
		MaxIdle:    MaxIdle,
		MaxCap:     MaxCap,
		Factory:    func() (interface{}, error) { return rpc.Dial(endpoint) },
		Close:      func(v interface{}) error { v.(*rpc.Client).Close(); return nil },
		// The maximum idle time of the connection, the connection exceeding this time will be closed,
		// which can avoid the problem of automatic failure when connecting to EOF when idle
		IdleTimeout: time.Duration(IdleTimeOut) * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		return
	}
	return EthRpcPool{Pool: p, opts: opts}, nil
}

// GetClient Get a connection from the connection pool
func (e *EthRpcPool) GetClient() (*rpc.Client, error) {
	client, err := e.Get()
	return client.(*rpc.Client), err
}

// ConnectLen View the number of connections in the current connection pool
func (e *EthRpcPool) ConnectLen() int {
	return e.Len()
}

// PutClient Put the connection back into the connection pool, when the connection is no longer in use
func (e *EthRpcPool) PutClient(client *rpc.Client) {
	_ = e.Put(client)
}

// ReleaseAll Release all connections in the connection pool, when resources need to be destroyed
func (e *EthRpcPool) ReleaseAll() {
	e.Release()
}

// BatchEthCall batch eth call
func (e *EthRpcPool) BatchEthCall(elements []rpc.BatchElem, allOk bool) (err error) {
	if len(elements) == 0 {
		return
	}

	shard := len(elements)/ChunkSize + 1

	var wg sync.WaitGroup
	errs := make([]error, shard)
	for i := 0; i < shard; i++ {
		wg.Add(1)

		start := i * ChunkSize
		end := (i + 1) * ChunkSize

		// 整除
		if start == len(elements) {
			wg.Done()
			break
		}
		if end > len(elements) {
			end = len(elements)
		}

		batch := elements[start:end]
		go func(batchIdx int, batch []rpc.BatchElem) {
			defer wg.Done()
			client, err := e.GetClient()
			for _, opt := range e.opts {
				opt(client)
			}
			if err != nil {
				return
			}
			defer e.PutClient(client)
			errs[batchIdx] = client.BatchCall(batch)
			if errs[batchIdx] != nil {
				log.Warn("batch call failed", "idx", batchIdx, "batchData", batch, "err", errs[batchIdx])
			}
		}(i, batch)
	}
	wg.Wait()
	for idx := range errs {
		if errs[idx] != nil {
			return fmt.Errorf("not all ok, occur a error, batchIdx is %d, err is %w", idx, errs[idx])
		}
	}
	if allOk {
		for idx := range elements {
			if elements[idx].Error != nil {
				return fmt.Errorf("not all ok, occur a error, elementIdx is %d, err is %w", idx, elements[idx].Error)
			}
		}
	}
	return
}

// Run 执行eth rpc
func (e *EthRpcPool) Run(f func(client *rpc.Client) error) error {
	client, err := e.GetClient()
	for _, opt := range e.opts {
		opt(client)
	}
	if err != nil {
		return err
	}
	defer e.PutClient(client)
	return f(client)
}
