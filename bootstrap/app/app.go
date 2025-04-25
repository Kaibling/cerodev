package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaibling/apiforge/ctxkeys"
	"github.com/kaibling/apiforge/log"
	apiservice "github.com/kaibling/apiforge/service"
	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/bootstrap/api"
	"github.com/kaibling/cerodev/config"
	"github.com/kaibling/cerodev/migration"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/repo/sqliterepo"
)

const volumePermissions = 0o755

func New() error {
	cfg := config.Load()
	baselogger := apiservice.BuildLogger(apiservice.LogConfig{ //nolint:exhaustruct
		LogLevel:     "debug",
		RequestBody:  true,
		ResponseBody: false,
		JSON:         false,
		AppName:      config.AppName,
	})
	appLogger := baselogger.Named("startup")
	ctx, ctxCancel := context.WithCancel(context.Background())

	conn, err := sqliterepo.Connect(cfg.DBConfig.FilePath)
	if err != nil {
		ctxCancel()

		return err
	}

	ctx = context.WithValue(ctx, ctxkeys.LoggerKey, baselogger)
	ctx = context.WithValue(ctx, ctxkeys.DBConnKey, conn)
	ctx = context.WithValue(ctx, ctxkeys.AppConfigKey, cfg)

	if err := migration.Migrate(conn, cfg); err != nil {
		appLogger.Warn("failed to migrate database: %s", err.Error())
		ctxCancel()

		return err
	}

	if err := ensureAdminUser(ctx); err != nil {
		appLogger.Warn("failed to ensure admin user: %s", err.Error())
		ctxCancel()

		return err
	}

	if err := ensurePorts(ctx); err != nil {
		appLogger.Warn("failed to ensure ports: %s", err.Error())
		ctxCancel()

		return err
	}

	if err := ensureVolumePath(cfg, appLogger); err != nil {
		ctxCancel()

		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if err := api.Start(ctx, cfg, baselogger, conn); err != nil {
		baselogger.Error("failed to start api", err)
	}

	appLogger.Info("application started. Ready...")
	<-interrupt

	appLogger.Info("stopping application")
	ctxCancel()

	gracePeriod := 400 * time.Millisecond //nolint:mnd
	time.Sleep(gracePeriod)

	return nil
}

func ensurePorts(ctx context.Context) error {
	l, ok := ctxkeys.GetValue(ctx, ctxkeys.LoggerKey).(log.Writer)
	if !ok {
		return errors.New("logger not found in context") //nolint:err113
	}

	cfg, ok := ctxkeys.GetValue(ctx, ctxkeys.AppConfigKey).(config.Configuration)
	if !ok {
		return errors.New("cfg not found in context") //nolint:err113
	}

	cs, err := bootstrap.NewContainerService(ctx)
	if err != nil {
		return err
	}

	pc, err := cs.GetPortCount()
	if err != nil {
		return err
	}

	if pc == 0 {
		l.Info("no ports found, creating new ones")

		if err := cs.FillPorts(cfg.ContainerMinPort, cfg.ContainerMaxPort); err != nil {
			return err
		}

		l.Info("ports created")
	}

	return nil
}

func ensureAdminUser(ctx context.Context) error {
	l, ok := ctxkeys.GetValue(ctx, ctxkeys.LoggerKey).(log.Writer)
	if !ok {
		return errors.New("logger not found in context") //nolint:err113
	}

	cfg, ok := ctxkeys.GetValue(ctx, ctxkeys.AppConfigKey).(config.Configuration)
	if !ok {
		return errors.New("cfg not found in context") //nolint:err113
	}

	adminUsername := "admin"

	us, err := bootstrap.NewUserService(ctx)
	if err != nil {
		return err
	}

	var adminUser *model.User

	adminUser, err = us.GetUnsafeByUsername(adminUsername)
	if err != nil {
		newAdminUser := &model.User{ //nolint:exhaustruct
			Username: adminUsername,
			Password: cfg.AdminPassword,
		}

		adminUser, err = us.Create(newAdminUser)
		if err != nil {
			l.Warn("failed to create admin user: %s", err.Error())

			return err
		}
	}

	if cfg.AdminToken != "" {
		return ensureToken(ctx, adminUser.ID, cfg.AdminToken)
	}

	return nil
}

func ensureToken(ctx context.Context, userID, token string) error {
	l, ok := ctxkeys.GetValue(ctx, ctxkeys.LoggerKey).(log.Writer)
	if !ok {
		return errors.New("logger not found in context") //nolint:err113
	}
	// check if admin token exists
	ts, err := bootstrap.NewTokenService(ctx)
	if err != nil {
		return err
	}

	savedToken, err := ts.GetByTokenKey(token)
	if err == nil {
		if savedToken.Token != token {
			// token is different, delete the old one
			err = ts.Delete(savedToken.Token)
			if err != nil {
				return err
			}

			l.Info("deleted old admin token")
		} else {
			return nil
		}
	}

	// create token
	newAdminToken := &model.Token{
		UserID: userID,
		Token:  token,
	}

	_, err = ts.CreateUnsafe(newAdminToken)
	if err != nil {
		return err
	}

	l.Info("created new admin token")
	l.Info("admin token rotated")

	return nil
}

func ensureVolumePath(cfg config.Configuration, l log.Writer) error {
	if _, err := os.Stat(cfg.VolumesPath); os.IsNotExist(err) {
		l.Info("volumes path does not exist, creating it")

		if err := os.MkdirAll(cfg.VolumesPath, volumePermissions); err != nil {
			l.Warn("failed to create volumes path: %s", err.Error())

			return err
		}

		l.Info("volumes path created")
	}

	return nil
}
