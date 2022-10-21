package main

import (
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pipeline"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/u"
	"runtime"
	"time"
)

func main() {
	if err := u.InitializeEnv("./.env"); err != nil {
		panic("Failed to load env file.")
	}

	if err := u.InitializeLogger(u.GetEnvWithDefault("ENVIRONMENT", "development")); err != nil {
		panic("Failed to initialize logger.")
	}

	switchboardConfig, err := config.LoadConfig("./config.yaml")
	if err != nil {
		u.Logger.Error(err.Error())
	}

	pubsubClient, err := pubsub.GetPubSubClient()
	if err != nil {
		u.Logger.Fatal(err.Error())
	}

	err = pubsubClient.Connect()
	if err != nil {
		u.Logger.Fatal(err.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(ginzap.Ginzap(u.Logger.ZapLogger, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(u.Logger.ZapLogger, true))

	for _, route := range switchboardConfig.Routes {
		server.GET(route.Endpoint, pipeline.NewRoutePipeline(route))
	}

	if u.GetEnvWithDefault("ENVIRONMENT", "development") == "development" {
		go func() {
			for {
				// TODO: Goroutine leak!
				u.Logger.Debug(fmt.Sprintf("%d running goroutines", runtime.NumGoroutine()))
				time.Sleep(1 * time.Second)
			}
		}()
	}

	err = server.Run(fmt.Sprintf("localhost:%d", switchboardConfig.Server.Port))
	if err != nil {
		u.Logger.Fatal("Failed to start gin server.")
	}

}
