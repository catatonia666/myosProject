package apiserver

import (
	"context"
	"dialogue/internal/models"
	"dialogue/internal/store/sqlstore"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	defer sqlDB.Close()

	store := sqlstore.New(db)

	redisDB, err := newCash(config.RedisURL)
	if err != nil {
		return err
	}

	s := newServer(store, redisDB)

	return http.ListenAndServe(config.BindAddr, s)
}

func newDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Migrator().DropTable("first_blocks", "blocks")
	db.AutoMigrate(&models.FirstBlock{}, &models.Block{}, &models.User{})

	return db, nil
}

func newCash(redisURL string) (*redis.Client, error) {
	var url *redis.Options
	url, _ = redis.ParseURL(redisURL)

	redisClient := redis.NewClient(url)

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	return redisClient, nil
}
