package pubsub

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jettdc/switchboard/pubsub/listen_groups"
	"github.com/jettdc/switchboard/u"
	"strconv"
	"time"
)

type RedisConnection struct {
	Client                    *redis.Client
	subscriptionListenHandler listen_groups.ListenGroupHandler
}

var Redis = &RedisConnection{nil, listen_groups.NewStdListenGroupHandler()}

// Connect establishes a connection to redis, using environment variables:
//   - REDIS_DATABASE_NUMBER (default 0)
//   - REDIS_ADDRESS (default "localhost:6379")
//   - REDIS_PASSWORD (default none)
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

// Subscribe leverages the base forwarder to only establish a single subscription to each topic. It also uses Redis'
// PSUBSCRIBE method to allow for parameterized topics.
func (r *RedisConnection) Subscribe(ctx context.Context, topic string, listenerId string) (chan listen_groups.ForwardedMessage, error) {
	return listen_groups.BaseForwarder(ctx, topic, r.subscriptionListenHandler, listenerId, redisSubscriptionRoutine)
}

func redisSubscriptionRoutine(topic string, doneChannel <-chan bool, messages chan<- listen_groups.ForwardedMessage, subscriptionDone chan<- bool, ctx context.Context) {
	pubsub := Redis.Client.PSubscribe(ctx, topic)
	defer pubsub.Close()

	for {
		select {
		case <-doneChannel:
			u.Logger.Debug(fmt.Sprintf("No more listeners on topic %s. Unsubscribing.", topic))
			subscriptionDone <- true
			return
		case msg := <-pubsub.Channel():
			packedMessage := listen_groups.ForwardedMessage{
				msg.Channel,
				msg.Payload,
			}

			messages <- packedMessage
		}
	}
}
