package appctx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaibling/apiforge/ctxkeys"
	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/config"
)

func GetBaseData(ctx context.Context) (*sql.DB, log.Writer, config.Configuration, error) { //nolint:ireturn
	db, ok := ctxkeys.GetValue(ctx, ctxkeys.DBConnKey).(*sql.DB)
	if !ok {
		return nil, nil, config.Configuration{}, errors.New("db connection not found in context") //nolint:err113
	}

	l, ok := ctxkeys.GetValue(ctx, ctxkeys.LoggerKey).(log.Writer)
	if !ok {
		return nil, nil, config.Configuration{}, errors.New("db connection not found in context") //nolint:err113
	}

	cfg, ok := ctxkeys.GetValue(ctx, ctxkeys.AppConfigKey).(config.Configuration)
	if !ok {
		return nil, nil, config.Configuration{}, errors.New("app config not found in context") //nolint:err113
	}

	return db, l, cfg, nil
}

func GetToken(ctx context.Context) (string, error) {
	token, ok := ctxkeys.GetValue(ctx, ctxkeys.TokenKey).(string)
	if !ok {
		return "", errors.New("token not found in context") //nolint:err113
	}

	return token, nil
}
