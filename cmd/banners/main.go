package main

import (
	"banners/internal/auth"
	"banners/internal/config"
	"banners/internal/handlers"
	"banners/internal/logger"
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase/BMS"
	"banners/internal/usecase/authentication"
	"banners/internal/usecase/cache"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	pxgConfig, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		logger.Fatalf("failed to parse config: %s", err)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, pxgConfig)
	if err != nil {
		logger.Fatalf("failed to connect to postgres: %s", err)
	}

	repo := banners.New(dbPool)

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Cache.RedisAddr})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Fatalf("failed to connect to redis: %s", err)
	}

	redisCache := cache.New(cache.Deps{
		RedisClient: redisClient,
	}, cfg.Cache.ExpirationTime)

	bms := BMS.NewBMS(BMS.Deps{
		Repository: repo,
		TxBuilder:  dbPool,
		Cache: redisCache,
	})

	authenticator := auth.New(cfg.Auth)

	authSystem := authentication.NewAuthenticationSystem(authentication.Deps{
		Authenticator: authenticator,
		Repo:          repo,
	})
	usecases := handlers.Usecases{
		BannerManagementSystem: bms,
		AuthenticationSystem:   authSystem,
	}
	controller := handlers.NewController(usecases)
	server := controller.NewServer(cfg.HTTPServer)
	log.Printf("server is listening at %s", cfg.HTTPServer.Address)
	log.Fatal(server.ListenAndServe())
}
