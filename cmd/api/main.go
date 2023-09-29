package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/cruffinoni/fizzbuzz/internal/api"
	"github.com/cruffinoni/fizzbuzz/internal/config"
	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/gin-gonic/gin"
)

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

	g := router.Group("/tasks")
	g.POST("/play", routes.PlayFizzBuzz)
	g.GET("/ping", routes.Ping)

	log.Printf("Starting api...")
	err = router.Run(fmt.Sprintf(":%d", configuration.APIPort))
	if err != nil {
		log.Fatalf("server ended with an error: %v", err)
	}
}
