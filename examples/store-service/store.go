package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Redis Client - TODO: fit to docker compose service
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

type Req struct {
	ID string `uri:"id" binding:"required"`
}

type Notification struct {
	Message string `json:"message"`
}

// Helpers
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
func storeStatus(c *gin.Context) {

	// Get ID from request
	var req Req
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Send Messages
	msg := new(Notification)
	msg.Message = "Order received"
	publishToRedis(*msg, req.ID)
	time.Sleep(5 * time.Second)

	msg.Message = "Pizza in oven"
	publishToRedis(*msg, req.ID)
	time.Sleep(10 * time.Second)

	msg.Message = "Pizza done cooking"
	publishToRedis(*msg, req.ID)
}

func main() {
	r := gin.Default()

	// Endpoints
	r.GET("/store/:id/events", storeStatus)

	r.Run(":54321") // listen and serve on localhost:12345
}
