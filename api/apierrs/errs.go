package apierrs

import (
	"errors"

	"github.com/kaibling/apiforge/apierror"
	"github.com/kaibling/cerodev/errs"
)

func HandleError(err error) apierror.HTTPError {
	if errors.Is(err, errs.ErrDataNotFound) {
		return apierror.ErrDataNotFound
	}

	return apierror.ErrServerError
}
