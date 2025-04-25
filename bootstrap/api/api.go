package api

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kaibling/apiforge/ctxkeys"
	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/apiforge/middleware"
	apiservice "github.com/kaibling/apiforge/service"
	"github.com/kaibling/apiforge/status"
	"github.com/kaibling/cerodev/api"
	"github.com/kaibling/cerodev/config"
	"github.com/kaibling/cerodev/web"
)

const CorsMaxAge = 300

func Start(
	ctx context.Context,
	cfg config.Configuration,
	baselogger log.Writer,
	conn *sql.DB,
) error {
	root := chi.NewRouter()
	// context
	root.Use(middleware.AddContext(ctxkeys.LoggerKey, baselogger))
	root.Use(middleware.AddContext(ctxkeys.DBConnKey, conn))
	root.Use(middleware.AddContext(ctxkeys.AppConfigKey, cfg))

	// middleware
	root.Use(cors.Handler(cors.Options{ //nolint:exhaustruct
		AllowedOrigins: []string{"*"},
		// Access-Control-Allow-Origin
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           CorsMaxAge,
	}))

	root.Use(middleware.InitEnvelope)
	root.Use(middleware.SaveBody)
	root.Use(middleware.LogRequest)
	root.Use(middleware.Recoverer)

	root.Mount("/api/v1", api.Route())
	web.AddUIRoute(root)

	// root.NotFound(handler.NotFound)

	apiServer := apiservice.New(ctx, apiservice.ServerConfig{ //nolint:exhaustruct
		BindingIP:   cfg.APIBindingIP,
		BindingPort: cfg.APIBindingPort,
	})

	status.IsReady.Store(true)

	return apiServer.Start(root, baselogger)
}
