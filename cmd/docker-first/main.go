package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"docker-first/internal/cache"
	"docker-first/internal/config"
	"docker-first/internal/database"
	"docker-first/internal/migrations"
	"docker-first/internal/service"

	"github.com/pkg/errors"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "error in config initiating"))
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error in create database conn"))
	}
	defer db.Close()

	err = migrations.Migrate(db, cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error in migrate process"))
	}

	redisCache, err := cache.InitCache(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error in cache initiating"))
	}

	cacheCtx, cancel := context.WithTimeout(context.Background(), cfg.CacheTimeout)
	defer cancel()
	err = redisCache.InitData(cacheCtx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error in init cache data"))
	}

	router := service.InitRouter(db, redisCache, cfg)

	log.Println("service starting...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println(errors.Wrap(err, "error in running service"))
		os.Exit(0)
	}
}
