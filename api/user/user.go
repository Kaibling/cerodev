package user

import (
	"net/http"

	"github.com/kaibling/apiforge/apierror"
	"github.com/kaibling/apiforge/envelope"
	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/errs"
	"github.com/kaibling/cerodev/errs/msg"
)

func usersGet(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_user")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	us, err := bootstrap.NewUserService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.UserServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	users, err := us.GetAll()
	if err != nil {
		l.Warn(errs.ErrMsg("cannot get all users", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(users).Finish(w, r, l)
}
