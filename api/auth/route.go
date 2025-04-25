package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/kaibling/cerodev/api/middleware"
)

func Route() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/login", login)
		r.With(middleware.Authentication).Group(func(r chi.Router) {
			r.Post("/logout", logout)
			r.Get("/check", check)
		})
	})

	return r
}
