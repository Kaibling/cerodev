package template

import (
	"github.com/go-chi/chi/v5"
	"github.com/kaibling/cerodev/api/middleware"
)

func Route() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Authentication)
		r.Post("/", createTemplate)
		r.Get("/", getTemplates)
		r.Delete("/{id}", deleteTemplate)
		r.Post("/{id}", buildImage)
		r.Put("/{id}", updateTemplate)
	})

	return r
}
