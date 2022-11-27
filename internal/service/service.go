package service

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"docker-first/internal/cache"
	"docker-first/internal/config"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Service interface {
	// routes
	GetCount(w http.ResponseWriter, r *http.Request)
	UpdateCount(w http.ResponseWriter, r *http.Request)
	About(w http.ResponseWriter, r *http.Request)
}

func NewService(db *sqlx.DB, redisCache cache.Cache, cfg *config.Config) Service {
	return &HttpService{
		db:         db,
		redisCache: redisCache,
		cfg:        cfg,
	}
}

type HttpService struct {
	db         *sqlx.DB
	redisCache cache.Cache
	cfg        *config.Config
}

func (s *HttpService) GetCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cacheCtx, cancel := context.WithTimeout(ctx, s.cfg.CacheTimeout)
	defer cancel()
	counterValue, err := s.redisCache.Get(cacheCtx)
	if err != nil {
		err = errors.Wrap(err, "error in get available currencies")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := `
		insert into counter_data (get_data, client_info, operation) values ($1, $2, $3)
	`

	dbCtx, cancel := context.WithTimeout(ctx, s.cfg.DBTimeout)
	defer cancel()
	res, err := s.db.ExecContext(dbCtx, query, time.Now(), r.Header.Get("User-Agent"), "get counter")
	if err != nil {
		err = errors.Wrap(err, "error in getting currency to ban data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		err = errors.New("error in insert userdata")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.FormatInt(counterValue, 10)))
	w.WriteHeader(http.StatusOK)
	return
}

func (s *HttpService) UpdateCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cacheCtx, cancel := context.WithTimeout(ctx, s.cfg.CacheTimeout)
	defer cancel()
	counterValue, err := s.redisCache.Increment(cacheCtx)
	if err != nil {
		err = errors.Wrap(err, "error in get available currencies")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := `
		insert into counter_data (get_data, client_info, operation) values ($1, $2, $3)
	`

	dbCtx, cancel := context.WithTimeout(ctx, s.cfg.DBTimeout)
	defer cancel()
	res, err := s.db.ExecContext(dbCtx, query, time.Now(), r.Header.Get("User-Agent"), "update counter")
	if err != nil {
		err = errors.Wrap(err, "error in getting currency to ban data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		err = errors.New("error in insert userdata")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(strconv.FormatInt(counterValue, 10)))
	w.WriteHeader(http.StatusOK)
	return
}

func (s *HttpService) About(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := `
		insert into counter_data (get_data, client_info, operation) values ($1, $2, $3)
	`

	dbCtx, cancel := context.WithTimeout(ctx, s.cfg.DBTimeout)
	defer cancel()
	res, err := s.db.ExecContext(dbCtx, query, time.Now(), r.Header.Get("User-Agent"), "about host")
	if err != nil {
		err = errors.Wrap(err, "error in getting currency to ban data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		err = errors.New("error in insert userdata")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name, ok := os.LookupEnv("NAME")
	if !ok {
		err = errors.New("no $NAME provided")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(name))
	w.WriteHeader(http.StatusOK)
	return
}
