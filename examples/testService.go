package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Redis Client
// TODO: extend Pub/Sub class and Redis implementation
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

type Test struct {
	Message string `json:"message"`
}

// GET /api/ws/test
// Triggers publishing a message to redis
func getTest(c *gin.Context) {

	var ctx = context.Background()
	testMsg := new(Test)

	payload, err := json.Marshal(testMsg)
	if err != nil {
		panic(err)
	}
	if err := redisClient.Publish(ctx, "/test/orders", payload).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Published message to redis")

}

func main() {
	r := gin.Default()

	// Endpoints
	r.GET("/api/ws/test", getTest)

	r.Run() // listen and serve on localhost:8080
}
