package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/kaibling/cerodev/api/auth"
	"github.com/kaibling/cerodev/api/container"
	images "github.com/kaibling/cerodev/api/image"
	"github.com/kaibling/cerodev/api/template"
	"github.com/kaibling/cerodev/api/user"
)

func Route() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Mount("/users", user.Route())
	r.Mount("/containers", container.Route())
	r.Mount("/templates", template.Route())
	r.Mount("/images", images.Route())
	r.Mount("/auth", auth.Route())

	return r
}
