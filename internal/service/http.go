package service

import (
	"net/http"

	"docker-first/internal/cache"
	"docker-first/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRouter(db *sqlx.DB, redisCache cache.Cache, cfg *config.Config) http.Handler {
	s := NewService(db, redisCache, cfg)

	r := chi.NewRouter()
	initMiddlewares(r)
	initRoutes(r, s)

	return r
}
