package middleware

import (
	"context"
	"net/http"
	"strings"

	apierror "github.com/kaibling/apiforge/apierror"
	"github.com/kaibling/apiforge/ctxkeys"
	"github.com/kaibling/apiforge/envelope"
	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/bootstrap"
	"github.com/kaibling/cerodev/errs"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read envelope
		e, l, aerr := envelope.GetEnvelopeAndLogger(r, "authentication")
		if aerr != nil {
			e.SetError(aerr).Finish(w, r, l)

			return
		}

		// read token
		tokenString, aerr := extractToken(r.Header, l)
		if aerr != nil {
			l.Warn("Token not extracted")
			e.SetError(apierror.ErrForbidden).Finish(w, r, l)

			return
		}

		// validate token and get username
		us, err := bootstrap.NewUserService(r.Context())
		if err != nil {
			l.Warn(errs.ServiceBuildError(bootstrap.UserServiceName, err))
			e.SetError(apierror.ErrForbidden).Finish(w, r, l)

			return
		}

		// todo set token last used
		user, err := us.CheckToken(tokenString)
		if err != nil {
			l.Warn("Error checking token: %s", err.Error())
			e.SetError(apierror.New(errs.ErrInvalidToken, http.StatusBadRequest)).Finish(w, r, l)

			return
		}

		ctx := context.WithValue(r.Context(), ctxkeys.UserNameKey, user.Username)
		ctx = context.WithValue(ctx, ctxkeys.UserIDKey, user.ID)
		ctx = context.WithValue(ctx, ctxkeys.TokenKey, tokenString)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(header http.Header, l log.Writer) (string, apierror.HTTPError) {
	// add logger and remove prints
	if _, ok := header["Authorization"]; !ok {
		l.Warn("Authorization header not found")

		return "", apierror.ErrForbidden
	}

	if len(header["Authorization"]) != 1 {
		l.Warn("Multiple Authorization tokens found")

		return "", apierror.ErrForbidden
	}

	authSlice := strings.Split(header["Authorization"][0], " ")

	position := 2
	if len(authSlice) != position {
		l.Warn("Bearer token format not valid")

		return "", apierror.ErrForbidden
	}

	return authSlice[1], nil
}
