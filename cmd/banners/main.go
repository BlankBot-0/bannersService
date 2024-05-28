package main

import (
	"banners/internal/auth"
	"banners/internal/config"
	"banners/internal/handlers"
	"banners/internal/repository/postgres/banners"
	"banners/internal/usecase/BMS"
	"banners/internal/usecase/authentication"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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
