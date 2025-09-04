package database

import (
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	pgDB        *gorm.DB
	redisClient *redis.Client
)

func InitDatabases(pgConfig PostgresConfig, redisConfig RedisConfig) error {
	var err error
	pgDB, err = NewPostgresConnection(pgConfig)
	if err != nil {
		return err
	}

	redisClient, err = NewRedisConnection(redisConfig)
	if err != nil {
		return err
	}
	return nil
}

func InitPostgres(pgConfig PostgresConfig) *gorm.DB {
	var err error
	counter := 0
	for {
		pgDB, err = NewPostgresConnection(pgConfig)
		if err != nil {
			log.Printf("error in initialize postgres %v", err)
		} else {
			return pgDB
		}
		if counter > 5 {
			return nil
		}
	}

}
func GetPostgres() *gorm.DB {
	return pgDB
}

func GetRedis() *redis.Client {
	return redisClient
}
