package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pipeline"
	"github.com/jettdc/switchboard/pubsub"
	"log"
)

func main() {
	switchboardConfig, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	pubsubClient, err := pubsub.GetPubSubClient()
	if err != nil {
		log.Fatal(err)
	}

	err = pubsubClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	server := gin.Default()

	for _, route := range switchboardConfig.Routes {
		server.GET(route.Endpoint, pipeline.NewRoutePipeline(route))
	}

	err = server.Run(fmt.Sprintf("localhost:%d", switchboardConfig.Server.Port))
	if err != nil {
		log.Fatal("Failed to start gin server.")
	}

}
