package pubsub

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jettdc/switchboard/u"
	"strconv"
	"time"
)

type RedisConnection struct {
	Client *redis.Client
}

var Redis = &RedisConnection{}

func (r *RedisConnection) Connect() error {
	u.Logger.Info("Connecting to Redis")

	// Already connected
	if Redis.Client != nil {
		return nil
	}

	dbEnv := u.GetEnvWithDefault("REDIS_DATABASE_NUMBER", "0")
	dbNo, err := strconv.Atoi(dbEnv)
	if err != nil {
		return fmt.Errorf("invalid redis database number: must be integer")
	}

	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     u.GetEnvWithDefault("REDIS_ADDRESS", "localhost:6379"),
		Password: u.GetEnvWithDefault("REDIS_PASSWORD", ""),
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

	u.Logger.Info("Successfully connected to Redis.")

	return nil
}

func (r *RedisConnection) Subscribe(ctx context.Context, topic string) (chan Message, error) {
	u.Logger.Debug(fmt.Sprintf("Subscribing to redis topic %s", topic))
	messages := make(chan Message, 8)

	go func() {
		pubsub := Redis.Client.PSubscribe(ctx, topic)
		
		defer pubsub.Close()
		defer u.Logger.Debug(fmt.Sprintf("Unsubscribing from redis topic %s", topic))

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-pubsub.Channel():
				packedMessage := Message{
					msg.Channel,
					msg.Payload,
				}

				messages <- packedMessage
			}
		}
	}()

	return messages, nil
}
