package dbrepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaibling/cerodev/errs"
)

func ToAppError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", errs.ErrDataNotFound, err) //nolint:errorlint
	}

	if errors.Is(err, sql.ErrTxDone) {
		return fmt.Errorf("%w: %v", errs.ErrDataTxError, err) //nolint:errorlint
	}

	return fmt.Errorf("%w: %v", errs.ErrInternalError, err) //nolint:errorlint
}
