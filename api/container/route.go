package container

import (
	"github.com/go-chi/chi/v5"
	"github.com/kaibling/cerodev/api/middleware"
)

func Route() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Authentication)
		r.Post("/", createContainer)
		r.Get("/", getContainers)
		r.Delete("/{id}", deleteContainer)
		r.Post("/{id}/start", startContainer)
		r.Post("/{id}/stop", stopContainer)
	})

	return r
}
