package auth

import (
	"net/http"

	"github.com/kaibling/apiforge/apierror"
	"github.com/kaibling/apiforge/envelope"
	"github.com/kaibling/apiforge/route"
	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/bootstrap/appctx"
	"github.com/kaibling/cerodev/errs"
	"github.com/kaibling/cerodev/errs/msg"
	"github.com/kaibling/cerodev/model"
)

func login(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_auth")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	var loginRequest model.LoginRequest
	if err := route.ReadPostData(r, &loginRequest); err != nil {
		l.Warn(errs.ErrMsg(msg.RequestParse, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	us, err := bootstrap.NewUserService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.UserServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	newToken, err := us.Login(&loginRequest)
	if err != nil {
		l.Warn(msg.RequestParse, err)
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(newToken).Finish(w, r, l)
}

func logout(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_auth")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	token, err := appctx.GetToken(r.Context())
	if err != nil {
		l.Warn(errs.ErrMsg("cannot get token", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	ts, err := bootstrap.NewTokenService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.UserServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	err = ts.Delete(token)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot delete token", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetSuccess().Finish(w, r, l)
}

func check(w http.ResponseWriter, r *http.Request) {
	e, l, merr := envelope.GetEnvelopeAndLogger(r, "api_auth")
	if merr != nil {
		l.Warn(errs.ErrMsg(msg.EnvelopeLoad, merr))
		e.SetError(merr).Finish(w, r, l)

		return
	}

	token, err := appctx.GetToken(r.Context())
	if err != nil {
		l.Warn(errs.ErrMsg("cannot get token", err))

		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	us, err := bootstrap.NewUserService(r.Context())
	if err != nil {
		l.Warn(errs.ServiceBuildError(bootstrap.UserServiceName, err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	user, err := us.CheckToken(token)
	if err != nil {
		l.Warn(errs.ErrMsg("cannot check token", err))
		e.SetError(apierror.NewGeneric(err)).Finish(w, r, l)

		return
	}

	e.SetResponse(user).Finish(w, r, l)
}
