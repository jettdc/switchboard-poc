package pubsub

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jettdc/switchboard/common"
	"log"
	"strconv"
	"time"
)

type RedisConnection struct {
	Client *redis.Client
}

var Redis = &RedisConnection{}

func (r *RedisConnection) Connect() error {
	log.Println("Connecting to Redis")

	// Already connected
	if Redis.Client != nil {
		return nil
	}

	dbEnv := common.GetEnvWithDefault("REDIS_DATABASE_NUMBER", "0")
	dbNo, err := strconv.Atoi(dbEnv)
	if err != nil {
		return fmt.Errorf("invalid redis database number: must be integer")
	}

	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     common.GetEnvWithDefault("REDIS_ADDRESS", "localhost:6379"),
		Password: common.GetEnvWithDefault("REDIS_PASSWORD", ""),
		DB:       dbNo,
	})

	_, err = Redis.Client.Ping(Redis.Client.Context()).Result()
	if err != nil {
		// Might just not be initialized
		time.Sleep(3)
		err := Redis.Client.Ping(Redis.Client.Context()).Err()
		if err != nil {
			return fmt.Errorf("failed to establish redis connection %s", err.Error())
		}
	}

	log.Println("Successfully connected to Redis.")

	return nil
}

func (r *RedisConnection) Subscribe(ctx context.Context, topic string) (chan string, error) {
	log.Printf("Subscribing to redis topic %s", topic)
	messages := make(chan string, 8)

	go func() {
		// TODO: Only psubscribe if *?
		pubsub := Redis.Client.PSubscribe(ctx, topic)
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-pubsub.Channel():
				packedMessage, err := PubSubMessage{
					msg.Channel,
					msg.Payload,
				}.String()
				if err != nil {
					log.Println(err)
					continue
				}

				// TODO: Enrich messages
				messages <- packedMessage
			}
		}

	}()

	return messages, nil
}
