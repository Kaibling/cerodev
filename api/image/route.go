package images

import (
	"github.com/go-chi/chi/v5"
	"github.com/kaibling/cerodev/api/middleware"
)

func Route() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Authentication)
		r.Get("/", getImages)
	})

	return r
}
