package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisMQ struct {
	rdb *redis.Client
	msg chan *redis.Message
}

var mq *RedisMQ

func NewMQ(ctx context.Context, cnn string) *RedisMQ {
	opt, err := redis.ParseURL(cnn)
	if err != nil {
		panic(fmt.Errorf("new mq failed, err is %v", err))
	}

	opt.ReadTimeout = -1 // set read time out to descord
	rdb := redis.NewClient(opt)

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("ping failed"))
	}

	mq = &RedisMQ{rdb: rdb, msg: make(chan *redis.Message)}
	return mq
}

func GetMQ() *RedisMQ {
	return mq
}

func (r *RedisMQ) Subscribe(ctx context.Context, channel string) {
	log.Info("debug mq", "mq conn is ok?", r.rdb.Ping(ctx))

	pubsub := r.rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Critical(ctx, err.Error())
		}
		mq.msg <- msg
	}
}

func (r *RedisMQ) Set(ctx context.Context, k string, v interface{}, t time.Duration) error {
	return r.rdb.Set(ctx, k, v, t).Err()
}

func (r *RedisMQ) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisMQ) Pop() *redis.Message {
	return <-mq.msg
}

func (r *RedisMQ) Publish(ctx context.Context, channel, payload string) error {
	err := r.rdb.Publish(ctx, channel, payload).Err()
	if err != nil {
		log.Critical(ctx, err.Error())
		return err
	}
	return nil
}
