package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	_ "github.com/cruffinoni/fizzbuzz/docs"
	"github.com/cruffinoni/fizzbuzz/internal/api"
	"github.com/cruffinoni/fizzbuzz/internal/config"
	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			FizzBuzz API
// @version		1.0
// @description	This is a customized FizzBuzz API server.
// @host			localhost:8080
// @BasePath		/
// @schemes		http
func main() {
	var configuration config.Global
	if err := env.Parse(&configuration, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatalf("error during parsing config: %v", err)
	}
	db, err := database.NewDB(&configuration.Database)
	if err != nil {
		log.Fatalf("can't initialize connection to the database: %v", err)
	}

	if configuration.Environment != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	routes := api.NewRoutes(db)
	router := gin.New()

	router.POST("/play", routes.PlayFizzBuzz)
	router.GET("/most-used", routes.GetMostUsedRequest)
	router.GET("/ping", routes.Ping)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Printf("Starting api...")
	err = router.Run(fmt.Sprintf(":%d", configuration.APIPort))
	if err != nil {
		log.Fatalf("server ended with an error: %v", err)
	}
}
