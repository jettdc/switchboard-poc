package main

import (
	"fmt"
	"net/http"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/jettdc/switchboard/config"
	"github.com/jettdc/switchboard/pipeline"
	"github.com/jettdc/switchboard/pubsub"
	"github.com/jettdc/switchboard/u"
	"github.com/jettdc/switchboard/login"
	"time"
)

func main() {
	if err := u.InitializeLogger(u.GetEnvWithDefault("ENVIRONMENT", "development")); err != nil {
		panic("Failed to initialize logger.")
	}

	switchboardConfig, err := config.LoadConfig("./config.yaml")
	if err != nil {
		u.Logger.Fatal(err.Error())
	}

	if switchboardConfig.Server.EnvFile != "" {
		if err := u.InitializeEnv(switchboardConfig.Server.EnvFile); err != nil {
			u.Logger.Fatal(fmt.Sprintf("Failed to load env file at path \"%s\"", switchboardConfig.Server.EnvFile))
		}
	}

	pubsubClient, err := pubsub.GetPubSubClient(switchboardConfig.Server.Pubsub.Provider)
	if err != nil {
		u.Logger.Fatal(err.Error())
	}

	err = pubsubClient.Connect()
	if err != nil {
		u.Logger.Fatal(err.Error())
	}

	//go func() {
	//	for {
	//		time.Sleep(2 * time.Second)
	//		u.Logger.Debug(fmt.Sprintf("%d", runtime.NumGoroutine()))
	//	}
	//}()

	u.Logger.Fatal(startServer(switchboardConfig, pubsubClient).Error())
}

func startServer(c *config.Config, pubsubClient pubsub.PubSub) error {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(ginzap.Ginzap(u.Logger.ZapLogger, time.RFC3339, true))
	server.Use(ginzap.RecoveryWithZap(u.Logger.ZapLogger, true))

	// Health check
	server.GET("/", func(c *gin.Context) { c.JSON(200, "OK") })

	// For requesting multiple logical websockets with a single connection
	server.GET("/multi", pipeline.NewDynamicRoutePipeline(c, pubsubClient))

	for _, route := range c.Routes {
		server.GET(route.Endpoint, pipeline.NewRoutePipeline(route, pubsubClient))
	}

	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	if c.Server.SSL == nil || c.Server.SSL.Mode == "none" {
		u.Logger.Info(fmt.Sprintf("Running server @ http://%s", addr))
		return server.Run()
	}

	switch c.Server.SSL.Mode {
	case "auto":
		u.Logger.Info(fmt.Sprintf("Running server @ https://%s", c.Server.Host))
		return autotls.Run(server, c.Server.Host)
	case "", "manual":
		u.Logger.Info(fmt.Sprintf("Running server @ https://%s", addr))
		return server.RunTLS(addr, c.Server.SSL.CertPath, c.Server.SSL.KeyPath)
	default:
		return fmt.Errorf("invalid ssl type")
	}

	//Login Auth
	router := gin.Default()
	router.POST("/loginJSON", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err == nil {
			fmt.Println("json receive - %+v", json.Username)
			
		}else{
			fmt.Println("error - %+v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": json,
		})

		fmt.Println("get username as: ", json.Username)

		foundPassword := login.PostgresQuery(json.Username)
		if foundPassword == json.Password{
			fmt.Println("password is the same; login success")
		}else{
			fmt.Println("login failed!")
		}
	})
	router.Run("localhost:8080")
}
