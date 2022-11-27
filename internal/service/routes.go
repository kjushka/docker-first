package service

import (
	"github.com/go-chi/chi/v5"
)

func initRoutes(r chi.Router, s Service) {
	r.Get("/", s.GetCount)
	r.Get("/stat", s.UpdateCount)
	r.Get("/about", s.About)
}
