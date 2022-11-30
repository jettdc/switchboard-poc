package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: getEnvVar("REDIS_ADDRESS", "localhost:6379"),
})

type Req struct {
	ID string `uri:"id" binding:"required"`
}

type Notification struct {
	Message string `json:"message"`
}

// Helpers
func getEnvVar(varName string, varDefault string) string {
	value, exists := os.LookupEnv(varName)
	if !exists {
		return varDefault
	}
	return value
}

func publishToRedis(msg Notification, id string) {
	// Marshal delivery message
	payload, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	// Publish to the correct topic
	var ctx = context.Background()
	topic := "/orders/" + id + "/events"
	if err := redisClient.Publish(ctx, topic, payload).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Published message to ", topic)
}

// GET /delivery/:id:events
// Triggers publishing a message to redis
func deliveryStatus(c *gin.Context) {

	// Get ID from request
	var req Req
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Send Messages
	msg := new(Notification)
	msg.Message = "Delivery Started"
	publishToRedis(*msg, req.ID)
	time.Sleep(5 * time.Second)

	msg.Message = "15 min away"
	publishToRedis(*msg, req.ID)
	time.Sleep(10 * time.Second)

	msg.Message = "5 min away"
	publishToRedis(*msg, req.ID)
	time.Sleep(5 * time.Second)

	msg.Message = "Delivered"
	publishToRedis(*msg, req.ID)
}

func main() {
	r := gin.Default()

	// Endpoints
	r.GET("/delivery/:id/events", deliveryStatus)

	r.Run(":12345") // listen and serve on localhost:12345
}
