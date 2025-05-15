package bootstrap

import (
	"context"
	"errors"

	"github.com/kaibling/apiforge/ctxkeys"
	"github.com/kaibling/cerodev/bootstrap/appctx"
	"github.com/kaibling/cerodev/pkg/docker"
	"github.com/kaibling/cerodev/pkg/repo/dbrepo"
	"github.com/kaibling/cerodev/pkg/ws"
	"github.com/kaibling/cerodev/service"
)

const (
	UserServiceName      string = "user_service"
	TokenServiceName     string = "token_service"
	ContainerServiceName string = "container_service"
	TemplateServiceName  string = "template_service"
)

func NewUserService(ctx context.Context) (*service.UserService, error) {
	db, l, cfg, err := appctx.GetBaseData(ctx)
	if err != nil {
		return nil, err
	}

	ur := dbrepo.NewUserRepo(ctx, db, l)
	tr := dbrepo.NewTokenRepo(ctx, db, l)
	ts := service.NewTokenService(tr, cfg)

	return service.NewUserService(ur, ts, cfg), nil
}

func NewContainerService(ctx context.Context) (*service.ContainerService, error) {
	db, l, cfg, err := appctx.GetBaseData(ctx)
	if err != nil {
		return nil, err
	}

	dr := docker.NewRepo(ctx, cfg.VolumesPath)
	cr := dbrepo.NewContainerRepo(ctx, db, l)
	tr := dbrepo.NewTemplateRepo(ctx, db, l)

	return service.NewContainerService(cr, dr, tr, l, cfg), nil
}

func NewTemplateService(ctx context.Context) (*service.TemplateService, error) {
	db, l, _, err := appctx.GetBaseData(ctx)
	if err != nil {
		return nil, err
	}

	ur := dbrepo.NewTemplateRepo(ctx, db, l)

	return service.NewTemplateService(ur), nil
}

func NewTokenService(ctx context.Context) (*service.TokenService, error) {
	db, l, cfg, err := appctx.GetBaseData(ctx)
	if err != nil {
		return nil, err
	}

	tr := dbrepo.NewTokenRepo(ctx, db, l)

	return service.NewTokenService(tr, cfg), nil
}

func GetWebSocketService(ctx context.Context) (*service.WebSocketService, error) {
	ws, ok := ctxkeys.GetValue(ctx, "websocket").(*service.WebSocketService)
	if !ok {
		return nil, errors.New("websocker service not found in context") //nolint:err113
	}
	return ws, nil
}

func NewWebSocketService() (*service.WebSocketService, error) {
	wsr := ws.New()

	return service.NewWebSocketService(wsr), nil
}
