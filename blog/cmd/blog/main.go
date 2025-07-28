package main

import (
	"go-blog/internal/database"
	"log"

	"go-blog/internal/config"
	"go-blog/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(cfg.Server.Port)

	// init dbs
	_ = database.InitDatabases(database.NewPostgresConfig(), database.RedisConfig(cfg.Redis))
	db := database.GetPostgres()
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB connection: %v", err)
	}

	defer sqlDb.Close()

	//Initialize Redis
	// redisClient := database.GetRedis()
	// defer redisClient.Close()

	r := gin.Default()

	r.LoadHTMLGlob("../../web/*")

	router := routes.SetRouter(db)

	// Start server
	router.Run(":8080")
}
