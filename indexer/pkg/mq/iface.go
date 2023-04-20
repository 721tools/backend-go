package mq

import (
	"context"

	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
)

var (
	log = log16.NewLogger("module", "mq")
)

type MQ interface {
	Subscribe(ctx context.Context, channel string)
	pop()
	Publish(ctx context.Context, channel, payload string) error
	NewMq(connection_url string)
	GetMq()
}
