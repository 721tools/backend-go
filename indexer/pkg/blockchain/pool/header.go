package pool

import "github.com/ethereum/go-ethereum/rpc"

type Option func(client *rpc.Client)

func WithHeaders(headers map[string]string) Option {
	return func(client *rpc.Client) {
		for key, val := range headers {
			client.SetHeader(key, val)
		}
	}
}
