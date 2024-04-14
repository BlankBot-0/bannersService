package main

import (
	"banners/internal/config"
	"banners/internal/handlers"
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase/BMS"
	"banners/internal/usecase/authentification"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jellydator/ttlcache/v3"
	"log"
	"time"
)

func main() {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](30 * time.Minute),
	)
	cache.Set()

	ctx := context.Background()
	cfg := config.MustLoad()
	// temporary line
	fmt.Println(cfg)
	pxgConfig, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, pxgConfig)
	if err != nil {
		log.Fatal(err)
	}

	repo := banners.New(dbPool)

	bms := BMS.NewBMS(BMS.Deps{
		Repository: repo,
		TxBuilder:  dbPool,
	})
	auth := authentification.NewAuthentificationSystem(authentification.Deps{
		Repo: repo,
	}, cfg.Auth)
	usecases := handlers.Usecases{
		BannerManagementSystem: bms,
		AuthentificationSystem: auth,
	}
	controller := handlers.NewController(usecases)
	server := controller.NewServer(cfg.HTTPServer)
	log.Printf("server is listening at %s", cfg.HTTPServer.Address)
	log.Fatal(server.ListenAndServe())
}
