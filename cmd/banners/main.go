package main

import (
	"banners/internal/config"
	"banners/internal/handlers"
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase/BMS"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

func main() {
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
	usecases := handlers.Usecases{BannerManagementSystem: bms}
	controller := handlers.NewController(usecases)
	router := controller.NewRouter()
	addr := ":8080"
	log.Printf("server is listening at %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
